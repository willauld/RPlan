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

	"github.com/fatih/structs"
	"github.com/spf13/pflag"
	//pflag "flag"

	"github.com/willauld/lpsimplex"
	"github.com/willauld/rplanlib"
)

var version = struct {
	major int
	minor int
	str   string
}{0, 3, "rc2"}

//__version__ = '0.3-rc2'

/* TODO add to unit testing
def precheck_consistancy():
    print("\nDoing Pre-check:")
    # check that there is income for all contibutions
        #tcontribs = 0
    for year in range(S.numyr):
        t = 0
        for j in range(len(S.accounttable)):
            if S.accounttable[j]['acctype'] != 'aftertax':
                v = S.accounttable[j]
                c = v.get('contributions', None)
                if c is not None:
                    t += c[year]
        if t > S.income[year]:
            print("year: %d, total contributions of (%.0f) to all Retirement accounts exceeds other earned income (%.0f)"%(year, t, S.income[year]))
            print("Please change the contributions in the toml file to be less than non-SS income.")
            exit(1)
    return True

def consistancy_check(res, years, taxbins, cgbins, accounts, accmap, vindx):
    # check to see if the ordinary tax brackets are filled in properly
    print()
    print()
    print("Consistancy Checking:")
    print()

    result = vvar.my_check_index_sequence(years, taxbins, cgbins, accounts, accmap, vindx)

    for year in range(S.numyr):
        s = 0
        fz = False
        fnf = False
        i_mul = S.i_rate ** (S.preplanyears+year)
        for k in range(len(taxinfo.taxtable)):
            cut, size, rate, base = taxinfo.taxtable[k]
            size *= i_mul
            s += res.x[vindx.x(year,k)]
            if fnf and res.x[vindx.x(year,k)] > 0:
                print("Inproper packed brackets in year %d, bracket %d not empty while previous bracket not full." % (year, k))
            if res.x[vindx.x(year,k)]+1 < size:
                fnf = True
            if fz and res.x[vindx.x(year,k)] > 0:
                print("Inproperly packed tax brackets in year %d bracket %d" % (year, k))
            if res.x[vindx.x(year,k)] == 0.0:
                fz = True
        if S.accmap['aftertax'] > 0:
            scg = 0
            fz = False
            fnf = False
            for l in range(len(taxinfo.capgainstable)):
                cut, size, rate = taxinfo.capgainstable[l]
                size *= i_mul
                bamount = res.x[vindx.y(year,l)]
                scg += bamount
                for k in range(len(taxinfo.taxtable)-1):
                    if taxinfo.taxtable[k][0] >= taxinfo.capgainstable[l][0] and taxinfo.taxtable[k][0] < taxinfo.capgainstable[l+1][0]:
                        bamount += res.x[vindx.x(year,k)]
                if fnf and bamount > 0:
                    print("Inproper packed CG brackets in year %d, bracket %d not empty while previous bracket not full." % (year, l))
                if bamount+1 < size:
                    fnf = True
                if fz and bamount > 0:
                    print("Inproperly packed GC tax brackets in year %d bracket %d" % (year, l))
                if bamount == 0.0:
                    fz = True
        TaxableOrdinary = OrdinaryTaxable(year)
        if (TaxableOrdinary + 0.1 < s) or (TaxableOrdinary - 0.1 > s):
            print("Error: Expected (age:%d) Taxable Ordinary income %6.2f doesn't match bracket sum %6.2f" %
                (year + S.startage, TaxableOrdinary,s))

        for j in range(len(S.accounttable)):
            a = res.x[vindx.b(year+1,j)] - (res.x[vindx.b(year,j)] - res.x[vindx.w(year,j)] + deposit_amount(S, res, year, j))*S.accounttable[j]['rate']
            if a > 1:
                v = S.accounttable[j]
                print("account[%d], type %s, index %d, mykey %s" % (j, v['acctype'], v['index'], v['mykey']))
                print("account[%d] year to year balance NOT OK years %d to %d" % (j, year, year+1))
                print("difference is", a)

        T,spendable,tax,rate,cg_tax,earlytax,rothearly = IncomeSummary(year)
        if spendable + 0.1 < res.x[vindx.s(year)]  or spendable -0.1 > res.x[vindx.s(year)]:
            print("Calc Spendable %6.2f should equal s(year:%d) %6.2f"% (spendable, year, res.x[vindx.s(year)]))
            for j in range(len(S.accounttable)):
                print("+w[%d,%d]: %6.0f" % (year, j, res.x[vindx.w(year,j)]))
                print("-D[%d,%d]: %6.0f" % (year, j, deposit_amount(S, res, year, j)))
            print("+o[%d]: %6.0f +SS[%d]: %6.0f -tax: %6.0f -cg_tax: %6.0f" % (year, S.income[year] ,year, S.SS[year] , tax ,cg_tax))

        bt = 0
        for k in range(len(taxinfo.taxtable)):
            bt += res.x[vindx.x(year,k)] * taxinfo.taxtable[k][2]
        if tax + 0.1 < bt  or tax -0.1 > bt:
            print("Calc tax %6.2f should equal brackettax(bt)[]: %6.2f" % (tax, bt))
    print()
*/

/*
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

func printInputParams(ip *rplanlib.InputParams) {
	fmt.Printf("InputParams:\n")
	m := structs.Map(ip)
	i := 0
	for k, v := range m {
		if v != "" {
			fmt.Printf("%3d::'%30s': '%#v'\n", i, k, v)
		}
		i++
	}
}

func printInputParamsStrMap(m map[string]string) {
	fmt.Printf("InputParamsStrMap:\n")
	//fmt.Printf("ip: map[string]string{\n")
	for i, v := range rplanlib.InputStrDefs {
		if m[v] != "" {
			fmt.Printf("%3d::'%30s': '%s'\n", i, v, m[v])
			//fmt.Printf("\"%s\": \"%s\",\n", v, m[v])
		}
	}
	for j := 1; j < rplanlib.MaxStreams+1; j++ {
		for i, v := range rplanlib.InputStreamStrDefs {
			lineno := i + len(rplanlib.InputStrDefs)
			k := fmt.Sprintf("%s%d", v, j)
			if m[k] != "" {
				fmt.Printf("%3d::'%30s': '%s'\n", lineno, k, m[k])
				//fmt.Printf("\"%s\": \"%s\",\n", k, m[k])
			}
		}
	}
	fmt.Printf("\n")
	//fmt.Printf("},\n")
}

func help() {
	fmt.Printf("\nUsage: %s [options]* configfile\n", filepath.Base(os.Args[0]))
	pflag.PrintDefaults()
	os.Exit(0)
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
		"Allow optomizer create deposits beyond those explicity specified")

	InputStrStrMapPtr := pflag.BoolP("inputstringmap", "M", false,
		"Display Input string map (key, value) for all input parameters")

	versionPtr := pflag.BoolP("version", "V", false,
		"Display the program version number and exit")

	helpPtr := pflag.BoolP("help", "h", false,
		"Displays this help message and exit")

	pflag.Parse()

	if *helpPtr {
		help()
	}

	if *versionPtr == true {
		//__version__ = '0.3-rc2'
		fmt.Printf("\t%s: Version %d.%d-%s\n", filepath.Base(os.Args[0]), version.major, version.minor, version.str)
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

	var tomlfile string
	tomlfile = pflag.Arg(0)

	ipsmp, err := getInputStringsMapFromToml(tomlfile)
	if err != nil {
		fmt.Printf("Error reading toml file: %s\n", err)
		os.Exit(0)
	}

	if *InputStrStrMapPtr {
		printInputParamsStrMap(*ipsmp)
	}

	ip, err := rplanlib.NewInputParams(*ipsmp, msgList)
	if err != nil {
		fmt.Printf("ARetirementPlanner: %s\n", err)
		os.Exit(1)
	}
	//printInputParams(ip)

	//fmt.Printf("InputParams: %#v\n", ip)
	//os.Exit(0)
	ti := rplanlib.NewTaxInfo(ip.FilingStatus)
	taxbins := len(*ti.Taxtable)
	cgbins := len(*ti.Capgainstable)
	vindx, err := rplanlib.NewVectorVarIndex(ip.Numyr, taxbins,
		cgbins, ip.Accmap, os.Stdout)
	if err != nil {
		fmt.Printf("ARetirementPlanner: %s\n", err)
		os.Exit(1)
	}

	logfile := os.Stdout
	if *logfilePtr != "" {
		logfile, err = os.Create(*logfilePtr)
		if err != nil {
			fmt.Printf("ARetirementPlanner: %s\n", err)
			os.Exit(1)
		}
	}

	csvfile := (*os.File)(nil)
	if *csvPtr != "" {
		csvfile, err = os.Create(*csvPtr)
		if err != nil {
			fmt.Printf("ARetirementPlanner: %s\n", err)
			os.Exit(1)
		}
	}

	RoundToOneK := true
	if *oneKPtr {
		RoundToOneK = false
	}

	// TODO looks like verbosePTR does nothing - investigate
	ms, err := rplanlib.NewModelSpecs(vindx, ti, *ip, *VerbosePtr,
		*depositsPtr, RoundToOneK, os.Stderr, logfile, csvfile, logfile, msgList)
	if err != nil {
		fmt.Printf("ARetirementPlanner: %s\n", err)
		os.Exit(1)
	}
	wc := msgList.GetWarningCount()
	if wc > 0 {
		fmt.Printf("%d Warning(s) found:\n", wc)
		for i := 0; i < wc; i++ {
			fmt.Printf("%s\n", msgList.GetWarning(i))
		}
	}
	msgList.ClearWarnings()

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
	disp := *VerbosePtr //true // false //true
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
		ms.PrintModelMatrix(c, a, b, notes, slack, bindingOnly)
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
			fmt.Printf("ARetirementPlanner: %s\n", err)
			os.Exit(1)
		}
		// rplanlib.BinCheckModelFiles("./RPlanModelgo.datX", "./RPlanModelpython.datX", &vindx)
	}

	//fmt.Printf("Res: %#v\n", res)
	if *VerbosePtr && false {
		str := fmt.Sprintf("Message: %v\n", res.Message)
		fmt.Printf(str)
		fmt.Printf("\n")
		fmt.Printf("Num Vars:        %d\n", len(a[0]))
		fmt.Printf("Num Constraints: %d\n", len(a))
		//fmt.Printf("Called LPSimplex() for m:%d x n:%d model\n", len(a), len(a[0]))
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
			ms.PrintCapGainsBrackets(&res.X)
		}
		ms.PrintBaseConfig(&res.X)
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
