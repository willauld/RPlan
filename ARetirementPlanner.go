package main

//
// A Retirement Planner (optimize withdrawals for most efficient use
// of the nest egg)
//

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/spf13/pflag"
	//pflag "flag"

	"github.com/willauld/lpsimplex"
	"github.com/willauld/rplanlib"
)

var version = struct {
	major int
	minor int
	patch int
	str   string
}{0, 4, 0, "rc1"}

//__version__ = '0.4.0-rc1'

// TODO add to unit testing

/*
// TODO FIXME This should be added to the library test suite
def verifyInputs( c , A , b ):
    m = len(A)
    n = len(A[0])
    if len(c) != n :
        print("lp: c vector incorrect length")
    if len(b) != m :
        print("lp: b vector incorrect length")

	# Do some sanity checks so that ab does not become singular during the
	# simplex solution. If the ZeroRow checks are removed then the code for
	# finding a set of linearly indepent columns must be improved.

	# Check that if a row of A only has zero elements that corresponding
	# element in b is zero, otherwise the problem is infeasible.
	# Otherwise return ErrZeroRow.
    zeroRows = 0
    for i in range(m):
        isZero = True
        for j in range(n) :
            if A[i][j] != 0 :
                isZero = False
                break
        if isZero and b[i] != 0 :
            # Infeasible
            print("ErrInfeasible -- row[%d]\n"% i)
        elif isZero :
            zeroRows+=1
            print("ErrZeroRow -- row[%d]\n"% i)
    # Check that if a column only has zero elements that the respective C vector
    # is positive (otherwise unbounded). Otherwise return ErrZeroColumn.
    zeroColumns = 0
    for j in range( n) :
        isZero = True
        for i in range( m) :
            if A[i][j] != 0 :
                isZero = False
                break
        if isZero and c[j] < 0 :
            print("ErrUnbounded -- column[%d] %s\n"% (j, vindx.varstr(j)))
        elif isZero :
            zeroColumns+=1
            print("ErrZeroColumn -- column[%d] %s\n"% (j, vindx.varstr(j)))
    print("\nZero Rows: %d, Zero Columns: %d\n"%(zeroRows, zeroColumns))
*/
/*
// This is not strickly needed Now that I am using pflag.lookup to set a default
// ToDo: think about removing,
func commandLineFlagWasSet(flag string) bool {
	setflags := make(map[string]bool)

	pflag.Visit(func(f *pflag.Flag) { setflags[f.Name] = true })

	if setflags[flag] {
		// command line set this flag
		fmt.Printf("command line set flag: %s\n", flag)
		return true
	}
	// command line did NOT set this flag
	fmt.Printf("command line did NOT set flag: %s\n", flag)
	return false
}
*/

func getInputStrStrMapFromFile(f string) (*map[string]string, error) {
	file, err := os.Open(f)
	if err != nil {
		e := fmt.Errorf("Error: %s", err)
		return nil, e
	}
	defer file.Close()

	ipsm := rplanlib.NewInputStringsMap()

	scanner := bufio.NewScanner(file)

	// Default scanner is bufio.ScanLines. Lets use ScanWords.
	// Could also use a custom function of SplitFunc type
	//scanner.Split(bufio.ScanWords)

	var key, val string
	// Scan for next token.
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Printf("Line: %s\n", line)
		tokens := strings.SplitAfter(line, "'")
		first := true
		havePair := false
		for _, token := range tokens {
			if strings.Index(token, ":") != -1 {
				continue
			}
			//fmt.Printf("token: %s\n", token)
			if first && strings.Index(token, "'") != -1 {
				key = strings.TrimRight(token, "'")
				key = strings.TrimSpace(key)
				//fmt.Printf("key: %s\n", key)
				first = false
			} else if strings.Index(token, "'") != -1 {
				val = strings.TrimRight(token, "'")
				val = strings.TrimSpace(val)
				//fmt.Printf("val: %s\n", val)
				havePair = true
				break
			}
		}
		if havePair {
			//fmt.Printf("key: %s, val: %s\n", key, val)
			err := setStringMapValueWithValue(&ipsm, key, val)
			if err != nil {
				e := fmt.Errorf("Error: %s", err)
				return nil, e
			}
		}
	}
	// False on error or EOF. Check error
	err = scanner.Err()
	if err != nil {
		e := fmt.Errorf("Error: %s", err)
		return nil, e
	}
	return &ipsm, nil
}

func printInputParams(ip *rplanlib.InputParams) {
	f := os.Stdout
	fmt.Fprintf(f, "InputParams:\n")
	m := structs.Map(ip)
	i := 0
	for k, v := range m {
		if v != "" {
			fmt.Fprintf(f, "%3d::'%30s': '%#v'\n", i, k, v)
		}
		i++
	}
}

// List all keys from input string string map
func listInputParamsStrMap(f *os.File) {
	fmt.Printf("InputParamsStrMap Keys:\n")
	for i, v := range rplanlib.InputStrDefs {
		//fmt.Printf("%3d::'%s'\n", i, v)
		fmt.Fprintf(f, "%3d::'%30s': ''\n", i, v)
	}
	for j := 1; j < rplanlib.MaxStreams+1; j++ {
		for i, v := range rplanlib.InputStreamStrDefs {
			lineno := i + len(rplanlib.InputStrDefs) +
				(j-1)*len(rplanlib.InputStreamStrDefs)
			k := fmt.Sprintf("%s%d", v, j)
			//fmt.Printf("%3d::'%s'\n", lineno, k)
			fmt.Fprintf(f, "%3d::'%30s': ''\n", lineno, k)
		}
	}
	fmt.Printf("\n")
}

// print out the active input string string map
func printInputParamsStrMap(m map[string]string, f *os.File) {
	fmt.Fprintf(f, "InputParamsStrMap:\n")
	//fmt.Printf("ip: map[string]string{\n")
	for i, v := range rplanlib.InputStrDefs {
		if m[v] != "" {
			fmt.Fprintf(f, "%3d::'%32s': '%s'\n", i, v, m[v])
			//fmt.Printf("\"%s\": \"%s\",\n", v, m[v])
		}
	}
	for j := 1; j < rplanlib.MaxStreams+1; j++ {
		for i, v := range rplanlib.InputStreamStrDefs {
			lineno := i + len(rplanlib.InputStrDefs) +
				(j-1)*len(rplanlib.InputStreamStrDefs)
			k := fmt.Sprintf("%s%d", v, j)
			if m[k] != "" {
				fmt.Fprintf(f, "%3d::'%32s': '%s'\n", lineno, k, m[k])
				//fmt.Printf("\"%s\": \"%s\",\n", k, m[k])
			}
		}
	}
	fmt.Fprintf(f, "\n")
	//fmt.Printf("},\n")
}

func help() {
	fmt.Printf("\nUsage: %s [options]* configfile\n", filepath.Base(os.Args[0]))
	pflag.PrintDefaults()
	os.Exit(0)
}

func printMsgAndExit(msgList *rplanlib.WarnErrorList, err error) {
	errstr := fmt.Sprintf("%s", err)
	ec := msgList.GetErrorCount()
	found := false
	if ec > 0 {
		for i := 0; i < ec; i++ {
			msg := msgList.GetError(i)
			if errstr == msg {
				found = true
			}
		}
		if !found {
			fmt.Printf("%s\n", errstr)
		}
	}
	printMsg(msgList)
	os.Exit(0)
}

func printMsg(msgList *rplanlib.WarnErrorList) {
	// FIXME TODO
	ec := msgList.GetErrorCount()
	if ec > 0 {
		fmt.Printf("%d Error(s) found:\n", ec)
		for i := 0; i < ec; i++ {
			fmt.Printf("%s\n", msgList.GetError(i))
		}
	}
	msgList.ClearErrors()
	//

	wc := msgList.GetWarningCount()
	if wc > 0 {
		fmt.Printf("%d Warning(s) found:\n", wc)
		for i := 0; i < wc; i++ {
			fmt.Printf("%s\n", msgList.GetWarning(i))
		}
	}
	msgList.ClearWarnings()
}

func main() {

	//parser = argparse.ArgumentParser(description='Create an optimized finacial plan for retirement.')
	VerbosePtr := pflag.BoolP("verbose", "v", false,
		"Extra output from solver")

	AllPlanTablesPtr := pflag.BoolP("allplantables", "A", false,
		"Display all plan tables, equivalent to -aitb")

	AccountTransPtr := pflag.BoolP("accountTrans", "a", false,
		"Display the account transaction table")

	IncomePtr := pflag.BoolP("income", "i", false,
		"Display the income and expense table")

	TaxPtr := pflag.BoolP("tax", "t", false,
		"Display the tax info table")

	TaxBracketPtr := pflag.BoolP("brackets", "b", false,
		"Display the tax bracket info tables")

	ModelBindingPtr := pflag.BoolP("modelbinding", "o", false,
		"Write out the binding constraints for the LP model")

	ModelAllPtr := pflag.BoolP("modelall", "m", false,
		"Write out all the constraints for the LP model")

	dumpBinaryPtr := pflag.StringP("dumpbinary", "D", "",
		"Write a binary copy of the LP model (c, A, b, x?) to file")
	pflag.Lookup("dumpbinary").NoOptDefVal = "./RPlanModelgo.datX"

	loadBinaryPtr := pflag.StringP("loadbinary", "L", "",
		"Load a binary copy of the LP model (c, A, b, x?) from file")
	pflag.Lookup("loadbinary").NoOptDefVal = "./RPlanModelgo.datX"

	logfilePtr := pflag.StringP("logfile", "l", "",
		"Write RPlan results to Log file")
	pflag.Lookup("logfile").NoOptDefVal = "./RPlan.log"

	timePtr := pflag.BoolP("timesimplex", "S", false,
		"Display the time used by the simplex solver")

	csvPtr := pflag.StringP("csv", "c", "",
		"Additionally write output to a csv file")
	pflag.Lookup("csv").NoOptDefVal = "./RPlan.csv"

	oneKPtr := pflag.BoolP("nokrounding", "k", false,
		"Do not round results output to thousands")

	depositsPtr := pflag.BoolP("allowdeposits", "z", false,
		"Allow optomizer to create deposits beyond those explicity specified")

	developerPtr := pflag.BoolP("developerinfo", "d", false,
		"Extra output information for development")

	taxYearPtr := pflag.IntP("taxyear", "Y", 2018,
		"Set the year for the tax code to be used (currently 2017 and 2018 only)")

	OutputStrStrMapPtr := pflag.StringP("outputstringmap", "M", "",
		"Output Input string map (key, value) for all current input parameters (*.strmap)")
	pflag.Lookup("outputstringmap").NoOptDefVal = "stdout"

	InputStrStrMapKeysPtr := pflag.StringP("inputstrmaptemplate", "K", "",
		"Display Input string map for all possible input parameters (generates template (*.strmap))")
	pflag.Lookup("inputstrmaptemplate").NoOptDefVal = "stdout"

	fourPercentRulePtr := pflag.BoolP("4PercentRule", "4", false,
		"Experimental: Override the 'Desired Income' with 4.5% of assets")

	SPVRulePtr := pflag.BoolP("SPVRule", "5", false,
		"Experimental: Extra output for to move toward including Statistical Present Value SPV based plan evaluation")

	versionPtr := pflag.BoolP("version", "V", false,
		"Display the program version number and exit")

	helpPtr := pflag.BoolP("help", "h", false,
		"Displays this help message and exit")

	pflag.Parse()

	var err error

	inputstrmapfile := (*os.File)(os.Stdout)
	if *InputStrStrMapKeysPtr != "" {
		if *InputStrStrMapKeysPtr != "stdout" {
			inputstrmapfile, err = os.Create(*InputStrStrMapKeysPtr)
			if err != nil {
				fmt.Printf("Retirement Optimizer: %s\n", err)
				os.Exit(1)
			}
		}
		listInputParamsStrMap(inputstrmapfile)
	}

	if *helpPtr {
		help()
	}

	if *versionPtr == true {
		//__version__ = '0.3.0-rc2'
		fmt.Printf("\t%s: Version %d.%d.%d-%s\n", filepath.Base(os.Args[0]), version.major, version.minor, version.patch, version.str)
		os.Exit(0)
	}

	if pflag.NArg() < 1 {
		fmt.Printf("Error: Missing configation file name\n")
		help()
	}
	if pflag.NArg() > 1 {
		fmt.Printf("Error: Too many arguments for configation file name (%s)\n", pflag.Args())
		help()
	}

	msgList := rplanlib.NewWarnErrorList()

	infile := pflag.Arg(0)
	var ipsmp *map[string]string

	// infile can be .toml or .strmap, Toml file is assumed
	if filepath.Ext(infile) == ".strmap" {
		ipsmp, err = getInputStrStrMapFromFile(infile)
	} else {
		ipsmp, err = getInputStringsMapFromToml(infile)
	}
	if err != nil {
		e := fmt.Errorf("reading file (%s): %s", infile, err)
		printMsgAndExit(msgList, e)
	}

	strmapfile := (*os.File)(os.Stdout)
	if *OutputStrStrMapPtr != "" {
		if *OutputStrStrMapPtr != "stdout" {
			strmapfile, err = os.Create(*OutputStrStrMapPtr)
			if err != nil {
				fmt.Printf("Retirement Optimizer: %s\n", err)
				os.Exit(1)
			}
		}
		printInputParamsStrMap(*ipsmp, strmapfile)
	}

	ip, err := rplanlib.NewInputParams(*ipsmp, msgList)
	if err != nil {
		printMsgAndExit(msgList, err)
	}
	//printInputParams(ip)

	//fmt.Printf("InputParams: %#v\n", ip)
	//os.Exit(0)
	if *taxYearPtr != 2017 && *taxYearPtr != 2018 {
		fmt.Printf("Retirement Optimizer: %s\n",
			"only the 2017 and 2018 tax code years are supported")
		os.Exit(1)
	}
	ti := rplanlib.NewTaxInfo(ip.FilingStatus, *taxYearPtr)
	taxbins := len(*ti.Taxtable)
	cgbins := len(*ti.Capgainstable)
	vindx, err := rplanlib.NewVectorVarIndex(ip.Numyr, taxbins,
		cgbins, ip.Accmap, os.Stdout)
	if err != nil {
		printMsgAndExit(msgList, err)
	}

	logfile := os.Stdout
	if *logfilePtr != "" {
		logfile, err = os.Create(*logfilePtr)
		if err != nil {
			fmt.Printf("Retirement Optimizer: %s\n", err)
			os.Exit(1)
		}
	}

	csvfile := (*os.File)(nil)
	if *csvPtr != "" {
		csvfile, err = os.Create(*csvPtr)
		if err != nil {
			fmt.Printf("Retirement Optimizer: %s\n", err)
			os.Exit(1)
		}
	}

	RoundToOneK := true
	if *oneKPtr {
		RoundToOneK = false
	}

	// TODO looks like verbosePTR does nothing - investigate
	ms, err := rplanlib.NewModelSpecs(vindx, ti, *ip, *depositsPtr,
		RoundToOneK, *fourPercentRulePtr,
		os.Stderr, logfile, csvfile, logfile, msgList)
	if err != nil {
		printMsgAndExit(msgList, err)
	}

	//if commandLineFlagWasSet("loadbinary")
	//fmt.Printf("ModelSpecs: %#v\n", ms)
	if *loadBinaryPtr != "" {
		// this not only would need to load a model but vindx... to work correctly. maybe should not be here ... think about it
	}
	//}

	c, a, b, notes := ms.BuildModel()

	tol := 1.0e-7

	bland := false
	maxiter := 4000

	callback := lpsimplex.Callbackfunc(nil)
	//callback := lpsimplex.LPSimplexVerboseCallback
	//callback := lpsimplex.LPSimplexTerseCallback
	disp := false //*VerbosePtr //true // false //true

	start := time.Now()
	res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
	elapsed := time.Since(start)

	if *ModelAllPtr || *ModelBindingPtr {
		slack := []float64(nil)
		if res.Success {
			slack = res.Slack
		}
		bindingOnly := true
		if *ModelAllPtr || !res.Success {
			bindingOnly = false
		}
		ms.PrintModelMatrix(c, a, b, notes, slack, bindingOnly, nil)
	}

	//if commandLineFlagWasSet("dumpbinary")
	if *dumpBinaryPtr != "" {
		vid := &[]int32{
			int32(ip.Numyr),
			int32(taxbins),
			int32(cgbins),
			int32(ip.Accmap[rplanlib.IRA]),
			int32(ip.Accmap[rplanlib.Roth]),
			int32(ip.Accmap[rplanlib.Aftertax]),
		}
		err = rplanlib.BinDumpModel(c, a, b, res.X, vid, *dumpBinaryPtr)
		if err != nil {
			printMsgAndExit(msgList, err)
		}
		// rplanlib.BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
	}

	// Print all Error and Warning Messages
	printMsg(msgList)

	//fmt.Printf("Res: %#v\n", res)

	if *VerbosePtr /*&& false*/ {
		fmt.Printf("Num Vars:        %d\n", len(a[0]))
		fmt.Printf("Num Constraints: %d\n", len(a))
		fmt.Printf("res.Success: %v\n", res.Success)
	}
	if *timePtr {
		str := fmt.Sprintf("\nTime: LPSimplex() took %s\n", elapsed)
		fmt.Printf(str)
	}
	if res.Success {
		ms.PrintActivitySummary(&res.X)
		if *IncomePtr || *AllPlanTablesPtr {
			ms.PrintIncomeExpenseDetails()
		}
		if *AccountTransPtr || *AllPlanTablesPtr {
			ms.PrintAccountTrans(&res.X)
		}
		if *TaxPtr || *AllPlanTablesPtr {
			ms.PrintTax(&res.X)
		}
		if *TaxBracketPtr || *AllPlanTablesPtr {
			ms.PrintTaxBrackets(&res.X)
			if *developerPtr {
				ms.PrintShadowTaxBrackets(&res.X)
			}
			ms.PrintCapGainsBrackets(&res.X)
		}
		ms.PrintBaseConfig(&res.X)
		if *SPVRulePtr {
			ms.PrintAccountWithdrawals(&res.X) // TESTING TESTING TESTING FIXME TODO
		}
	} else {
		str := fmt.Sprintf("LP Simplex Message: %v\n", res.Message)
		fmt.Printf(str)
		fmt.Printf("\n")
		if ip.Min > 0 && ip.Maximize == rplanlib.PlusEstate {
			fmt.Printf("A possible cause when using [min.income] can be that the amount is more than can be supported by the assets over the plan period. Try lowering the amount if this may be the problem.")
		}
	}
	//createDefX(&res.X)
	//=-=-=-=-

	/* tODO Move this to NewInputParams
	   if S.accmap['IRA']+S.accmap['roth']+S.accmap['aftertax'] == 0:
	       print('Error: This app optimizes the withdrawals from your retirement account(s); you must have at least one specified in the input toml file.')
	       exit(0)
	*/
	/*
	   if precheck_consistancy():

	       #verifyInputs( c , A , b )
	       res = scipy.optimize.linprog(c, A_ub=A, b_ub=b,
	                                options={"disp": args.verbose,
	                                         #"bland": True,
	                                         "tol": 1.0e-7,
	                                         "maxiter": 4000})
	       consistancy_check(res, years, taxbins, cgbins, accounts, S.accmap, vindx)
	*/
}
