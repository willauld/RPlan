package main

//
// A Retirement Planner (optimize withdrawals for most efficient use
// of the nest egg)
//

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/pflag"

	"github.com/willauld/lpsimplex"
	"github.com/willauld/rplanlib"
)

var (
	// See createRelease.ps1 for variable definition / values
	vermajor         string
	verminor         string
	verpatch         string
	verstr           string
	buildTime        string
	gitLibHash       string
	gitDriverHash    string
	gitlpsimplexHash string
	// LAO is Local App Output var
	LAO rplanlib.AppOutput
)

var version = struct {
	major            string
	minor            string
	patch            string
	str              string
	buildTime        string
	gitLibHash       string
	gitDriverHash    string
	gitlpsimplexHash string
}{vermajor, verminor, verpatch,
	verstr, buildTime, gitLibHash, gitDriverHash, gitlpsimplexHash}

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

// List all possible keys from input string string map
// Can be used to create a template file
func listInputParamsStrMap(f *os.File) {
	fmt.Printf("InputParamsStrMap Keys:\n")
	for i, v := range rplanlib.InputStrDefs {
		fmt.Fprintf(f, "%3d::'%32s': ''\n", i, v)
	}
	for j := 1; j < rplanlib.MaxStreams+1; j++ {
		for i, v := range rplanlib.InputStreamStrDefs {
			lineno := i + len(rplanlib.InputStrDefs) +
				(j-1)*len(rplanlib.InputStreamStrDefs)
			k := fmt.Sprintf("%s%d", v, j)
			//fmt.Printf("%3d::'%s'\n", lineno, k)
			fmt.Fprintf(f, "%3d::'%32s': ''\n", lineno, k)
		}
	}
	fmt.Printf("\n")
}

func help() {
	fmt.Printf("\n'%s", filepath.Base(os.Args[0]))
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf(" %s", os.Args[i])
	}
	fmt.Printf("'\n\n")
	fmt.Printf("\nUsage: %s [options]* configfile\n", filepath.Base(os.Args[0]))
	pflag.PrintDefaults()
	os.Exit(0)
}

// GetrplanlibVersionString returns a version string
func GetrplanlibVersionString() string {
	v := rplanlib.Version
	s := fmt.Sprintf("%s.%s.%s", v.Major, v.Minor, v.Patch)
	if v.Str != "" {
		s = fmt.Sprintf("%s-%s", s, v.Str)
	}
	return s
}

// GetlpsimplexVersionString returns the lpsimplex library version string
func GetlpsimplexVersionString() string {
	v := lpsimplex.Version
	s := fmt.Sprintf("%s.%s.%s", v.Major, v.Minor, v.Patch)
	if v.Str != "" {
		s = fmt.Sprintf("%s-%s", s, v.Str)
	}
	return s
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
	}
	if !found {
		LAO.Output(fmt.Sprintf("%s\n", errstr))
	}
	printMsg(msgList)
	os.Exit(0)
}

func printMsg(msgList *rplanlib.WarnErrorList) {
	// FIXME TODO
	ec := msgList.GetErrorCount()
	if ec > 0 {
		LAO.Output(fmt.Sprintf("%d Error(s) found:\n", ec))
		for i := 0; i < ec; i++ {
			LAO.Output(fmt.Sprintf("%s\n", msgList.GetError(i)))
		}
	}
	msgList.ClearErrors()
	//

	wc := msgList.GetWarningCount()
	if wc > 0 {
		LAO.Output(fmt.Sprintf("%d Warning(s) found:\n", wc))
		for i := 0; i < wc; i++ {
			LAO.Output(fmt.Sprintf("%s\n", msgList.GetWarning(i)))
		}
	}
	msgList.ClearWarnings()
}

func main() {
	LAO = rplanlib.NewAppOutput(nil, nil) //defaults to stdout
	writeCmdLine := false

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
		"Write results to Log file")
	pflag.Lookup("logfile").NoOptDefVal = "./RPlan.log"

	timePtr := pflag.BoolP("timesimplex", "S", false,
		"Display the time used by the simplex solver")

	csvPtr := pflag.StringP("csv", "c", "",
		"Write output to a csv file")
	pflag.Lookup("csv").NoOptDefVal = "./RPlan.csv"

	oneKPtr := pflag.BoolP("nokrounding", "k", false,
		"Do not round results to thousands")

	developerPtr := pflag.BoolP("developerinfo", "d", false,
		"Extra output information for development")

	taxYearPtr := pflag.IntP("taxyear", "Y", 2018,
		"Set the year for the tax code to be used (currently 2017 and 2018 only)")

	OutputStrStrMapPtr := pflag.StringP("outputstringmap", "M", "",
		"Output string map (key, value) for all current input parameters (*.strmap)")
	pflag.Lookup("outputstringmap").NoOptDefVal = "stdout"

	InputStrStrMapKeysPtr := pflag.StringP("inputstrmaptemplate", "K", "",
		"Display string map for all possible input parameters (generates template (*.strmap))")
	pflag.Lookup("inputstrmaptemplate").NoOptDefVal = "stdout"

	LPoutputPtr := pflag.StringP("lpm", "f", "",
		"Write model out as LP_solve model format to a .lp file")
	pflag.Lookup("lpm").NoOptDefVal = "./RPlan.lp"

	// FIXME only works with piecewiseTaxfunction branch
	//ModelSpecsSetUsePiecewisePtr := pflag.BoolP("piecewise", "p", false,
	//	"Use piecewise linear tax calculation")

	fourPercentRulePtr := pflag.BoolP("4PercentRule", "4", false,
		"Experimental: Override the 'Desired Income' with 4.5% of assets")

	SPVRulePtr := pflag.BoolP("SPVRule", "5", false,
		"Experimental: Extra output for to move toward including Statistical Present Value SPV based plan evaluation")

	DynamicBlandPtr := pflag.BoolP("DynamicBland", "B", false,
		"Enable Bland Pivot Rule after each degenterate pivot")

	ScaleModelPtr := pflag.BoolP("ScaleModel", "E", false,
		"Equilibrate the model through scaling")

	versionPtr := pflag.BoolP("version", "V", false,
		"Display the program version number and exit")

	helpPtr := pflag.BoolP("help", "h", false,
		"Displays this help message and exit")

	pflag.Parse()

	var err error

	// writes an empty template with all possible input keys
	inputstrmapfile := (*os.File)(os.Stdout)
	if *InputStrStrMapKeysPtr != "" {
		if *InputStrStrMapKeysPtr != "stdout" {
			inputstrmapfile, err = os.Create(*InputStrStrMapKeysPtr)
			if err != nil {
				fmt.Printf("%s: %s\n", filepath.Base(os.Args[0]), err)
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
		if version.major == "" {
			str := GetrplanlibVersionString()
			str2 := GetlpsimplexVersionString()
			fmt.Printf("\t%s: No Release Version,\n\t\trplanlib: %s\n\t\tlpsimplex: %s\n",
				filepath.Base(os.Args[0]), str, str2)
		} else {
			fmt.Printf("\t%s: Version %s.%s.%s", filepath.Base(os.Args[0]), version.major, version.minor, version.patch)

			if version.str != "" {
				fmt.Printf("-%s\n", version.str)
			} else {
				fmt.Printf("\n")
			}

			if *VerbosePtr {
				fmt.Printf("\t\tBuild time:           %s\n", version.buildTime)
				fmt.Printf("\t\tDriver main Git Hash: %s\n", version.gitDriverHash)
				fmt.Printf("\t\trplanlib Git Hash:    %s\n", version.gitLibHash)
				fmt.Printf("\t\tlpsimplex Git Hash:    %s\n", version.gitlpsimplexHash)
			}
		}
		os.Exit(0)
	}

	if pflag.NArg() < 1 {
		fmt.Printf("\n'%s", filepath.Base(os.Args[0]))
		for i := 1; i < len(os.Args); i++ {
			fmt.Printf(" %s", os.Args[i])
		}
		fmt.Printf("'\n")
		fmt.Printf("Error: Missing configation file name\n")
		help()
	}
	if pflag.NArg() > 1 {
		fmt.Printf("\n'%s", filepath.Base(os.Args[0]))
		for i := 1; i < len(os.Args); i++ {
			fmt.Printf(" %s", os.Args[i])
		}
		fmt.Printf("'\n")
		fmt.Printf("Error: Too many arguments for configation file name (%s)\n", pflag.Args())
		help()
	}

	msgList := rplanlib.NewWarnErrorList()

	infile := pflag.Arg(0)
	var ipsmp *map[string]string

	// infile can be .toml or .strmap, Toml file is assumed
	if filepath.Ext(infile) == ".strmap" {
		ipsmp, err = rplanlib.GetInputStrStrMapFromFile(infile)
	} else {
		ipsmp, err = rplanlib.GetInputStringsMapFromToml(infile)
	}
	if err != nil {
		e := fmt.Errorf("reading file (%s): %s", infile, err)
		printMsgAndExit(msgList, e)
	}

	logfile := os.Stdout
	if *logfilePtr != "" {
		logfile, err = os.Create(*logfilePtr)
		if err != nil {
			fmt.Printf("%s: %s\n", filepath.Base(os.Args[0]), err)
			os.Exit(1)
		}
		writeCmdLine = true
		// TODO FIXME should a log file always start with the input parameters writen to the log? (Add this there?)
	}
	csvfile := (*os.File)(nil)
	if *csvPtr != "" {
		csvfile, err = os.Create(*csvPtr)
		if err != nil {
			fmt.Printf("%s: %s\n", filepath.Base(os.Args[0]), err)
			os.Exit(1)
		}
		writeCmdLine = true
		// TODO FIXME should a csv file always start with the input parameters writen to the csv? (Add this there?)
	}

	strmapfile := logfile
	if *OutputStrStrMapPtr != "" {
		if *OutputStrStrMapPtr != "stdout" {
			strmapfile, err = os.Create(*OutputStrStrMapPtr)
			if err != nil {
				fmt.Printf("%s: %s\n", filepath.Base(os.Args[0]), err)
				os.Exit(1)
			}
		}
		rplanlib.WriteFileInputParamsStrMap(strmapfile, *ipsmp)
	}

	ip, err := rplanlib.NewInputParams(*ipsmp, msgList)
	if err != nil {
		printMsgAndExit(msgList, err)
	}
	if *fourPercentRulePtr && ip.Maximize == rplanlib.PlusEstate {
		str := "Four percent rule and spending PlusEstate can not be used together"
		fmt.Printf("%s: %s\n", filepath.Base(os.Args[0]), str)
		os.Exit(1)
	}

	if *taxYearPtr != 2017 && *taxYearPtr != 2018 {
		fmt.Printf("%s: %s\n", filepath.Base(os.Args[0]),
			"only the 2017 and 2018 tax code years are supported")
		os.Exit(1)
	}
	LAO = rplanlib.NewAppOutput(csvfile, logfile)
	if writeCmdLine {
		LAO.Output(fmt.Sprintln(os.Args))
	}

	ti := rplanlib.NewTaxInfo(ip.FilingStatus, *taxYearPtr)
	taxbins := len(*ti.Taxtable)
	cgbins := len(*ti.Capgainstable)
	vindx, err := rplanlib.NewVectorVarIndex(ip.Numyr, taxbins,
		cgbins, ip.Accmap, os.Stdout)
	if err != nil {
		printMsgAndExit(msgList, err)
	}

	RoundToOneK := true
	if *oneKPtr {
		RoundToOneK = false
	}

	// TODO looks like verbosePTR does nothing - investigate
	ms, err := rplanlib.NewModelSpecs(vindx, ti, *ip,
		RoundToOneK, *developerPtr, *fourPercentRulePtr,
		os.Stderr, logfile, csvfile, logfile, msgList)
	if err != nil {
		printMsgAndExit(msgList, err)
	}
	// FIXME only works with piecewiseTaxfunction branch
	//if *ModelSpecsSetUsePiecewisePtr {
	//	ms.UsePieceWiseMethod = true
	//}

	//if commandLineFlagWasSet("loadbinary")
	//fmt.Printf("ModelSpecs: %#v\n", ms)
	if *loadBinaryPtr != "" {
		// this not only would need to load a model but vindx... to work correctly. maybe should not be here ... think about it
	}
	//}

	c, a, b, notes := ms.BuildModel()

	tol := 1.0e-7
	bland := false  // true   //false
	maxiter := 4000 // 4000

	callback := lpsimplex.Callbackfunc(nil)
	//callback := lpsimplex.LPSimplexVerboseCallback
	//callback := lpsimplex.LPSimplexTerseCallback
	disp := false //*VerbosePtr //true // false //true
	nb_cmd := lpsimplex.NB_CMD_RESET
	if *DynamicBlandPtr {
		nb_cmd = nb_cmd | lpsimplex.NB_CMD_USEDYNAMICBLAND
	}
	if *ScaleModelPtr {
		nb_cmd = nb_cmd | lpsimplex.NB_CMD_SCALEME | lpsimplex.NB_CMD_SCALEME_POW2
	}
	lpsimplex.LPSimplexSetNewBehavior(nb_cmd)

	start := time.Now()
	res := lpsimplex.LPSimplex(c, a, b, nil, nil, nil, callback, disp, maxiter, tol, bland)
	elapsed := time.Since(start)

	if *ModelAllPtr || *ModelBindingPtr {
		slack := []float64(nil)
		if res.Success || res.Status == 1 /*that is exceeded maxiter*/ {
			slack = res.Slack
		}
		bindingOnly := true
		if *ModelAllPtr || !res.Success {
			bindingOnly = false
		}
		//fmt.Printf("*** c[ms.Vindx.B(32,0)]: %6.4f\n", c[ms.Vindx.B(32, 0)])
		ms.PrintModelMatrix(c, a, b, notes, slack, bindingOnly, nil)
		ms.WriteObjectFunctionSolution(c, res.X, nil)
	}
	if *LPoutputPtr != "" {
		fmt.Printf("TESTING: *LPoutputPtr: %s\n", *LPoutputPtr)
		//LPoutputPtr := pflag.StringP("lpm", "L", "",
		// Write model in format consumable by lp_solve
		//fmt.Printf("*** c[ms.Vindx.B(32,0)]: %6.4f\n", c[ms.Vindx.B(32, 0)])
		ms.WriteLPFormatModel(c, a, b, notes, *LPoutputPtr, res.X, fmt.Sprintln(os.Args))
		//fmt.Printf("*** c[ms.Vindx.B(32,0)]: %6.4f\n", c[ms.Vindx.B(32, 0)])
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
		degenCount := lpsimplex.LPSimplexSetNewBehavior(lpsimplex.NB_CMD_NOP)
		fmt.Fprintf(logfile, "Num Vars:         %d\n", len(a[0]))
		fmt.Fprintf(logfile, "Num Constraints:  %d\n", len(a))
		fmt.Fprintf(logfile, "Iterations:       %d\n", res.Nitr)
		fmt.Fprintf(logfile, "Degen Iterations: %d\n", degenCount)
		fmt.Fprintf(logfile, "Objective func:   %15.4f\n", res.Fun)
		fmt.Fprintf(logfile, "res.Success:      %v\n", res.Success)
	}
	if *timePtr {
		str := fmt.Sprintf("\nTime: LPSimplex() took %s\n", elapsed)
		fmt.Fprintf(logfile, str)
	}
	if res.Success || (res.Status == 1 && *developerPtr) {
		ms.ConsistencyCheckBrackets(&res.X) // TODO FIXME this will be changed after debug
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
			ms.PrintAssetSummary()
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
