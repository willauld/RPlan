package main

//import "github.com/pelletier/go-toml"
import (
	"fmt"
	"reflect"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/willauld/rplanlib"
)

// Iam is the declearation of the user
type Iam struct {
	Primary                 bool
	Age                     int
	Retire                  int
	Through                 int
	DefinedContributionPlan string
}

var iamKeys = []string{"primary", "age", "retire", "through", "definedContributionPlan"}
var iamReqKeys = []string{"age", "retire", "through"}

// SocialSecurity is the Social Security information of the user
type SocialSecurity struct {
	Amount int
	Age    string
}

var socialSecurityKeys = []string{"amount", "age"}
var socialSecurityReqKeys = socialSecurityKeys

// Income is an income stream for the user
type Income struct {
	Amount    int
	Age       string
	Inflation bool
	Tax       bool
}

var incomeKeys = []string{"amount", "age", "inflation", "tax"}
var incomeReqKeys = []string{"amount", "age"}

// Expense if an expense stream for of the user
type Expense struct {
	Amount    int
	Age       string // FIXME TODO change toml def to be nicer eg this should be period rather than age
	Inflation bool
}

var expenseKeys = []string{"amount", "age", "inflation"}
var expenseReqKeys = []string{"amount", "age"}

// Asset is a user owned asset the may be sold
type Asset struct {
	Value               int
	CostAndImprovements int
	AgeToSell           int
	OwedAtAgeToSell     int
	PrimaryResidence    bool
	Rate                float64
	BrokerageRate       float64
}

var assetKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell", "primaryResidence", "rate", "brokerageRate"}
var assetReqKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell"}

// Desired is the minium after tax spending amount the user requires
type Desired struct {
	Amount int
}

var desiredKeys = []string{"amount"}
var desiredReqKeys = desiredKeys

// Max is the maximum after tax spending amount the user want to spend
type Max struct {
	Amount int
}

var maxKeys = []string{"amount"}
var maxReqKeys = maxKeys

// IRA is the IRA account info for the user
type IRA struct {
	Bal       int
	Rate      float64
	Contrib   int
	Inflation bool
	Period    string
}

var iRAKeys = []string{"bal", "rate", "contrib", "inflation", "period"}
var iRAReqKeys = []string{"bal"}

// Roth is the roth account info for the user
type Roth struct {
	Bal       int
	Rate      float64
	Contrib   int
	Inflation bool
	Period    string
}

var rothKeys = iRAKeys
var rothReqKeys = iRAReqKeys

// Aftertax is the aftertax account info for the user
type Aftertax struct {
	Bal       int
	Basis     int
	Rate      float64
	Contrib   int
	Inflation bool
	Period    string
}

var aftertaxKeys = []string{"bal", "basis", "rate", "contrib", "inflation", "period"}
var aftertaxReqKeys = iRAReqKeys

// GlobalConfig is all the information the user need enter for the retirement plan
type GlobalConfig struct {
	Title          string
	RetirementType string `toml:"retirement_type"`
	Returns        float64
	Inflation      float64
	Maximize       string
	Iam            map[string]Iam
	SocialSecurity map[string]SocialSecurity
	IRA            map[string]IRA
	Roth           map[string]Roth
	Aftertax       Aftertax
	Asset          map[string]Asset
	Income         map[string]Income
	Expense        map[string]Expense
}

var globalKeys = []string{"title", "retirement_type", "returns", "inflation", "maximize"}
var globalReqKeys = []string{"retiremen_type"}

var categoryKeys = []string{ //Top level key with global keys replace by single global
	"global",
	"iam",
	"SocialSecurity",
	"IRA",
	"roth",
	"aftertax",
	"asset",
	"income",
	"expense",
}

// TomlStrDefs works with InputStrDefs to Map Toml information to
// rplanlib API static portion
//
// List of toml paths to user supplied information.
// Includes editing codes for correcting as needed at runtime.
// '@' means the corresponding key in InputStrDefs should be used to set
//		iam key names
// '%n' means the nth iam key name is to be substituted here
// '%i' means the current stream key name is to be substituted here
// '#n' means that the resulting string is a range string and the nth
//		number from the range should be extracted and used for assignment
// After converting the toml paths there should not be any of these codes
// remaining.
// The resulting path and the corresponding key from InputStrDefs are used
// together to retreave and set the value in the Input String map.
var TomlStrDefs = []string{ // each entry corresponds to InputStrDefs entries
	"title",
	"retirement_type",
	"@iam.key1",
	"@iam.key2",
	"iam.%0.age",
	"iam.%1.age",
	"iam.%0.retire",
	"iam.%1.retire",
	"iam.%0.through",
	"iam.%1.through",
	"SocialSecurity.%0.amount",
	"SocialSecurity.%1.amount",
	"SocialSecurity.%0.age#1", //Need to extract from range
	"SocialSecurity.%1.age#1",
	"IRA.%0.bal",
	"IRA.%1.bal",
	"IRA.%0.rate",
	"IRA.%1.rate",
	"IRA.%0.contrib",
	"IRA.%1.contrib",
	"IRA.%0.period#1", //ContribStartAge, // need to extract from range #n means extract the number in the 1st or second posision in the result string "n1-n2"
	"IRA.%1.period#1", //ContribStartAge,
	"IRA.%0.period#2", //ContribEndAge,
	"IRA.%1.period#2", //ContribEndAge,
	"roth.%0.bal",
	"roth.%1.bal",
	"roth.%0.rate",
	"roth.%1.rate",
	"roth.%0.contrib",
	"roth.%1.contrib",
	"roth.%0.period#1", //contribStartAge, // need to extract from range
	"roth.%1.period#1", //contribStartAge,
	"roth.%0.period#2", //contribEndAge,
	"roth.%1.period#2", // contribEndAge,
	"aftertax.bal",
	"aftertax.rate",
	"aftertax.contrib",
	"aftertax.period#1", //contribStartAge, // need to extract from range
	"aftertax.period#2", //contribEndAge,

	"inflation",
	"returns",
	"maximize",
}

// TomlStreamStrDefs works with InputStreamStrDefs to Map Toml information to
// rplanlib API dynamic portion (per stream portion)
var TomlStreamStrDefs = []string{ // each entry corresponds to InputStreamStrDefs entries
	"@income",
	"income.%i.amount",
	"income.%i.age#1",
	"income.%i.age#2",
	"income.%i.inflation",
	"income.%i.tax",
	"@expense",
	"expense.%i.amount",
	"expense.%i.age#1",
	"expense.%i.age#2",
	"expense.%i.inflation",
	"expense.%i.tax",
	"@asset",
	"asset.%i.value",
	"asset.%i.ageToSell",
	"asset.%i.costAndImprovements",
	"asset.%i.owedAtAgeToSell",
	"asset.%i.primaryResidence",
	"asset.%i.rate",
	"asset.%i.brokerageRate",
}

func doItWithUnMarshal() {
	gc := GlobalConfig{}
	//toml.Unmarshal(document, &person)
	config, err := toml.LoadFile("hack.toml")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	config.Unmarshal(&gc)
	//TODO TESTME fmt.Printf("\ndoc22: %#v\n", gc)
}

// Required to have iam, retirement_type and at least one of IRA, roth or Aftertax

func keyIn(k string, keys []string) bool {
	for _, key := range keys {
		if k == key {
			return true
		}
	}
	return false
}

func categoryMatch(parent string, keys []string) bool {
	var reqKeys []string
	switch parent {
	case "retirement_type": // TODO FIXME change retirement_type to filingStatus
		fallthrough
	case "returns":
		fallthrough
	case "inflation":
		fallthrough
	case "maximize":
		return false
	case "iam":
		reqKeys = iamReqKeys
	case "SocialSecurity":
		reqKeys = socialSecurityReqKeys
	case "IRA":
		reqKeys = iRAReqKeys
	case "roth":
		reqKeys = rothReqKeys
	case "aftertax":
		reqKeys = aftertaxReqKeys
	case "income":
		reqKeys = incomeReqKeys
	case "asset":
		reqKeys = assetReqKeys
	case "expense":
		reqKeys = expenseReqKeys
	}
	for _, k := range reqKeys {
		if !keyIn(k, keys) {
			return false
		}
	}
	return true
}

func categoryMatchAndUnknowns(parent string, keys []string) (bool, []string) {
	var reqKeys []string
	var akeys []string
	unknownPaths := []string{}
	switch parent {
	case "retirement_type": // TODO FIXME change retirement_type to filingStatus
		fallthrough
	case "returns":
		fallthrough
	case "inflation":
		fallthrough
	case "maximize":
		return false, unknownPaths
	case "iam":
		reqKeys = iamReqKeys
		akeys = iamKeys
	case "SocialSecurity":
		reqKeys = socialSecurityReqKeys
		akeys = socialSecurityKeys
	case "IRA":
		reqKeys = iRAReqKeys
		akeys = iRAKeys
	case "roth":
		reqKeys = rothReqKeys
		akeys = rothKeys
	case "aftertax":
		reqKeys = aftertaxReqKeys
		akeys = aftertaxKeys
	case "income":
		reqKeys = incomeReqKeys
		akeys = incomeKeys
	case "asset":
		reqKeys = assetReqKeys
		akeys = assetKeys
	case "expense":
		reqKeys = expenseReqKeys
		akeys = expenseKeys
	}
	for _, k := range reqKeys {
		if !keyIn(k, keys) {
			return false, unknownPaths
		}
	}
	for _, k := range keys {
		if !keyIn(k, akeys) {
			//uk := parent + "." + k
			unknownPaths = append(unknownPaths, k)
		}
	}
	return true, unknownPaths
}

func getFloat64Value(obj interface{}) float64 {
	var targetVal float64
	switch reflect.TypeOf(obj).Name() {
	case "int64":
		targetVal = float64(obj.(int64))
	case "float64":
		targetVal = obj.(float64)
	}
	//fmt.Printf("getFloat64Value, got: %g\n", targetVal)
	return targetVal
}

func getInt64Value(obj interface{}) int64 {
	var targetVal int64
	switch reflect.TypeOf(obj).Name() {
	case "int64":
		targetVal = obj.(int64)
	case "float64":
		targetVal = int64(obj.(float64))
	}
	//fmt.Printf("getInt64Value, got: %d\n", targetVal)
	return targetVal
}

func getTomlGlobal(config *toml.Tree) GlobalConfig {
	globals := GlobalConfig{}
	for _, lPath := range globalKeys {
		lPathobj := config.Get(lPath)
		switch lPath {
		case "title":
			if config.Has(lPath) {
				globals.Title = lPathobj.(string)
				continue
			}
			globals.Title = ""
		case "retirement_type":
			if config.Has(lPath) {
				globals.RetirementType = lPathobj.(string)
				continue
			}
			globals.RetirementType = ""
		case "returns":
			if config.Has(lPath) {
				globals.Returns = getFloat64Value(lPathobj)
				continue
			}
			globals.Returns = -1
		case "inflation":
			if config.Has(lPath) {
				globals.Inflation = getFloat64Value(lPathobj)
				continue
			}
			globals.Inflation = -1
		case "maximize":
			if config.Has(lPath) {
				globals.Maximize = lPathobj.(string)
				continue
			}
			globals.Maximize = ""
		default:
			fmt.Printf("*** Implement ME(%s) ****\n", lPath)
		}
	}
	return globals
}

func getTomlAftertax(config *toml.Tree, path string) *Aftertax {

	if config.Has(path) {
		aftertaxCfig := Aftertax{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		//fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range aftertaxKeys {
				lPath := path + "." + v
				//fmt.Printf("**** lpath: %s ****\n", lPath)
				lPathobj := config.Get(lPath)
				switch v {
				case "bal":
					if config.Has(lPath) {
						aftertaxCfig.Bal = int(getInt64Value(lPathobj))
						continue
					}
					aftertaxCfig.Bal = -1
				case "basis":
					if config.Has(lPath) {
						aftertaxCfig.Basis = int(getInt64Value(lPathobj))
						continue
					}
					aftertaxCfig.Basis = -1
				case "contrib":
					if config.Has(lPath) {
						aftertaxCfig.Contrib = int(getInt64Value(lPathobj))
						continue
					}
					aftertaxCfig.Contrib = -1
				case "rate":
					if config.Has(lPath) {
						aftertaxCfig.Rate = getFloat64Value(lPathobj)
						continue
					}
					aftertaxCfig.Rate = -1
				case "inflation":
					if config.Has(lPath) {
						aftertaxCfig.Inflation = lPathobj.(bool)
						continue
					}
					aftertaxCfig.Inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "period":
					if config.Has(lPath) {
						aftertaxCfig.Period = lPathobj.(string)
						continue
					}
					aftertaxCfig.Period = ""
				}
			}
		} else {
			fmt.Printf("*** EERROOOORR ***\n")
			return nil
		}
		//fmt.Printf("aftertax config: %#v\n", aftertaxCfig)
		return &aftertaxCfig
	}
	return nil
}

func getTomlIAm(config *toml.Tree, path string) *Iam {

	if config.Has(path) {
		iamCfig := Iam{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		//fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range iamKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				switch v {
				case "age":
					if config.Has(lPath) {
						iamCfig.Age = int(getInt64Value(lPathobj))
						continue
					}
					iamCfig.Age = -1
				case "retire":
					if config.Has(lPath) {
						iamCfig.Retire = int(getInt64Value(lPathobj))
						continue
					}
					iamCfig.Retire = -1
				case "through":
					if config.Has(lPath) {
						iamCfig.Through = int(getInt64Value(lPathobj))
						continue
					}
					iamCfig.Through = -1
				case "definedContributionPlan":
					if config.Has(lPath) {
						iamCfig.DefinedContributionPlan = lPathobj.(string)
						continue
					}
					iamCfig.DefinedContributionPlan = ""
				case "primary":
					if config.Has(lPath) {
						iamCfig.Primary = lPathobj.(bool)
						continue
					}
					iamCfig.Primary = false
				}
			}
		} else {
			fmt.Printf("*** ???? ***\n")
			return nil
		}
		return &iamCfig
	}
	return nil
}

func getTomlIAmMap(config *toml.Tree, path string) *map[string]Iam {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	//fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
	m := make(map[string]Iam)
	if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
		m["nokey"] = *getTomlIAm(config, path)
		return &m
	}
	for _, k := range keys /*user define names*/ {
		lPath := path + "." + k
		//fmt.Printf("lPath: %s\n", lPath)
		//lPathT := config.Get(lPath).(*toml.Tree)
		//fmt.Printf("\nlpath: %s tree keys: %#v\n", lPath, lPathT.Keys())
		m[k] = *getTomlIAm(config, lPath)
	}
	return &m
}

func getTomlSocialSecurity(config *toml.Tree, path string) *SocialSecurity {

	if config.Has(path) {
		ssCfig := SocialSecurity{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			//var socialSecurityKeys = []string{"amount", "age"}
			for _, v := range socialSecurityKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				switch v {
				case "amount":
					if config.Has(lPath) {
						ssCfig.Amount = int(getInt64Value(lPathobj))
						continue
					}
					ssCfig.Amount = -1
				case "age":
					if config.Has(lPath) {
						ssCfig.Age = lPathobj.(string)
						continue
					}
					ssCfig.Age = ""
				}
			}
		} else {
			fmt.Printf("*** ???? ***\n")
			return nil
		}
		return &ssCfig
	}
	return nil
}

func getTomlSocialSecurityMap(config *toml.Tree, path string) *map[string]SocialSecurity {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	//fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
	m := make(map[string]SocialSecurity)
	if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
		m["nokey"] = *getTomlSocialSecurity(config, path)
		return &m
	}
	for _, k := range keys /*user define names*/ {
		lPath := path + "." + k
		//fmt.Printf("lPath: %s\n", lPath)
		//lPathT := config.Get(lPath).(*toml.Tree)
		//fmt.Printf("\nlpath: %s tree keys: %#v\n", lPath, lPathT.Keys())
		m[k] = *getTomlSocialSecurity(config, lPath)
	}
	return &m
}

func getTomlIRA(config *toml.Tree, path string) *IRA {

	if config.Has(path) {
		iraCfig := IRA{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range iRAKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				//var IRAKeys = []string{"bal", "rate", "contrib", "inflation", "period"}
				switch v {
				case "bal":
					if config.Has(lPath) {
						iraCfig.Bal = int(getInt64Value(lPathobj))
						continue
					}
					iraCfig.Bal = -1
				case "rate":
					if config.Has(lPath) {
						iraCfig.Rate = getFloat64Value(lPathobj)
						continue
					}
					iraCfig.Rate = -1
				case "contrib":
					if config.Has(lPath) {
						iraCfig.Contrib = int(getInt64Value(lPathobj))
						continue
					}
					iraCfig.Contrib = -1
				case "inflation":
					if config.Has(lPath) {
						iraCfig.Inflation = lPathobj.(bool)
						continue
					}
					iraCfig.Inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "period":
					if config.Has(lPath) {
						iraCfig.Period = lPathobj.(string)
						continue
					}
					iraCfig.Period = ""
				}
			}
		} else {
			fmt.Printf("*** ???? ***\n")
			return nil
		}
		return &iraCfig
	}
	return nil
}

func getTomlIRAMap(config *toml.Tree, path string) *map[string]IRA {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]IRA)
	if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
		m["nokey"] = *getTomlIRA(config, path)
		return &m
	}
	for _, k := range keys /*user define names*/ {
		lPath := path + "." + k
		m[k] = *getTomlIRA(config, lPath)
	}
	return &m
}

func getTomlRoth(config *toml.Tree, path string) *Roth {

	if config.Has(path) {
		rothCfig := Roth{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range rothKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				switch v {
				case "bal":
					if config.Has(lPath) {
						rothCfig.Bal = int(getInt64Value(lPathobj))
						continue
					}
					rothCfig.Bal = -1
				case "rate":
					if config.Has(lPath) {
						rothCfig.Rate = getFloat64Value(lPathobj)
						continue
					}
					rothCfig.Rate = -1
				case "contrib":
					if config.Has(lPath) {
						rothCfig.Contrib = int(getInt64Value(lPathobj))
						continue
					}
					rothCfig.Contrib = -1
				case "inflation":
					if config.Has(lPath) {
						rothCfig.Inflation = lPathobj.(bool)
						continue
					}
					rothCfig.Inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "period":
					if config.Has(lPath) {
						rothCfig.Period = lPathobj.(string)
						continue
					}
					rothCfig.Period = ""
				}
			}
		} else {
			fmt.Printf("*** ???? ***\n")
			return nil
		}
		return &rothCfig
	}
	return nil
}

func getTomlRothMap(config *toml.Tree, path string) *map[string]Roth {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]Roth)
	if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
		m["nokey"] = *getTomlRoth(config, path)
		return &m
	}
	for _, k := range keys /*user define names*/ {
		lPath := path + "." + k
		m[k] = *getTomlRoth(config, lPath)
	}
	return &m
}

func getTomlIncome(config *toml.Tree, path string) *Income {

	if config.Has(path) {
		iCfig := Income{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range incomeKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				//var incomeKeys = []string{"amount", "age", "inflation", "tax"}
				switch v {
				case "amount":
					if config.Has(lPath) {
						iCfig.Amount = int(getInt64Value(lPathobj))
						continue
					}
					iCfig.Amount = -1
				case "inflation":
					if config.Has(lPath) {
						iCfig.Inflation = lPathobj.(bool)
						continue
					}
					iCfig.Inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "tax":
					if config.Has(lPath) {
						iCfig.Tax = lPathobj.(bool)
						continue
					}
					iCfig.Tax = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "age":
					if config.Has(lPath) {
						iCfig.Age = lPathobj.(string)
						continue
					}
					iCfig.Age = ""
				}
			}
		} else {
			fmt.Printf("*** ???? ***\n")
			return nil
		}
		return &iCfig
	}
	return nil
}

func getTomlIncomeMap(config *toml.Tree, path string) *map[string]Income {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]Income)
	if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
		m["nokey"] = *getTomlIncome(config, path)
		return &m
	}
	for _, k := range keys /*user define names*/ {
		lPath := path + "." + k
		m[k] = *getTomlIncome(config, lPath)
	}
	return &m
}

func getTomlExpense(config *toml.Tree, path string) *Expense {

	if config.Has(path) {
		eCfig := Expense{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range expenseKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				//var expenseKeys = []string{"amount", "age", "inflation"}
				switch v {
				case "amount":
					if config.Has(lPath) {
						eCfig.Amount = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.Amount = -1
				case "inflation":
					if config.Has(lPath) {
						eCfig.Inflation = lPathobj.(bool)
						continue
					}
					eCfig.Inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "age":
					if config.Has(lPath) {
						eCfig.Age = lPathobj.(string)
						continue
					}
					eCfig.Age = ""
				}
			}
		} else {
			fmt.Printf("*** ???? ***\n")
			return nil
		}
		return &eCfig
	}
	return nil
}

func getTomlExpenseMap(config *toml.Tree, path string) *map[string]Expense {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]Expense)
	if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
		m["nokey"] = *getTomlExpense(config, path)
		return &m
	}
	for _, k := range keys /*user define names*/ {
		lPath := path + "." + k
		m[k] = *getTomlExpense(config, lPath)
	}
	return &m
}

func getTomlAsset(config *toml.Tree, path string) *Asset {

	if config.Has(path) {
		eCfig := Asset{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range assetKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				//var assetKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell", "primaryResidence", "rate", "broderageRate"}
				switch v {
				case "value":
					if config.Has(lPath) {
						eCfig.Value = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.Value = -1
				case "primaryResidence":
					if config.Has(lPath) {
						eCfig.PrimaryResidence = lPathobj.(bool)
						continue
					}
					eCfig.PrimaryResidence = false
					//TODO FIXME setting default bool to false IS THIS WHAT IT SHOULD BE???
				case "costAndImprovements":
					if config.Has(lPath) {
						eCfig.CostAndImprovements = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.CostAndImprovements = -1
				case "ageToSell":
					if config.Has(lPath) {
						eCfig.AgeToSell = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.AgeToSell = -1
				case "owedAtAgeToSell":
					if config.Has(lPath) {
						eCfig.OwedAtAgeToSell = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.OwedAtAgeToSell = -1
				case "rate":
					if config.Has(lPath) {
						eCfig.Rate = getFloat64Value(lPathobj)
						continue
					}
					eCfig.Rate = -1
				//var assetKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell", "primaryResidence", "rate", "broderageRate"}
				case "brokerageRate":
					if config.Has(lPath) {
						eCfig.BrokerageRate = getFloat64Value(lPathobj)
						continue
					}
					eCfig.BrokerageRate = -1
				}
			}
		} else {
			fmt.Printf("*** ???? ***\n")
			return nil
		}
		return &eCfig
	}
	return nil
}

func getTomlAssetMap(config *toml.Tree, path string) *map[string]Asset {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]Asset)
	if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
		m["nokey"] = *getTomlAsset(config, path)
		return &m
	}
	for _, k := range keys /*user define names*/ {
		lPath := path + "." + k
		m[k] = *getTomlAsset(config, lPath)
	}
	return &m
}

func goGetTomlData(filename string) (*GlobalConfig, *toml.Tree, error) {

	config, err := toml.LoadFile(filename)
	if err != nil {
		e := fmt.Errorf("Error: %s", err)
		return nil, nil, e
	}
	//fmt.Printf("\nKeys: %#v\n", config.Keys())
	//for _, b := range config.Keys()
	var gc GlobalConfig
	for _, path := range categoryKeys {
		switch path {
		case "global":
			gc = getTomlGlobal(config)
			//fmt.Printf("global configs: %#v\n", gc)
		case "iam":
			c := getTomlIAmMap(config, path)
			gc.Iam = *c
			//fmt.Printf("iam configs map: %#v\n", c)
		case "SocialSecurity":
			c := getTomlSocialSecurityMap(config, path)
			gc.SocialSecurity = *c
			//fmt.Printf("socail security configs map: %#v\n", c)
		case "IRA":
			c := getTomlIRAMap(config, path)
			gc.IRA = *c
			//fmt.Printf("IRA configs map: %#v\n", c)
		case "roth":
			c := getTomlRothMap(config, path)
			gc.Roth = *c
			//fmt.Printf("Roth configs map: %#v\n", c)
		case "aftertax":
			c := getTomlAftertax(config, path)
			gc.Aftertax = *c
			//fmt.Printf("aftertax configs: %#v\n", c)
		case "income":
			c := getTomlIncomeMap(config, path)
			gc.Income = *c
			//fmt.Printf("Income configs map: %#v\n", c)
		case "asset":
			c := getTomlAssetMap(config, path)
			gc.Asset = *c
			//fmt.Printf("Asset configs map: %#v\n", c)
		case "expense":
			c := getTomlExpenseMap(config, path)
			gc.Expense = *c
			//fmt.Printf("Expense configs map: %#v\n", c)
		default:
			pathT := config.Get(path).(*toml.Tree) //TomlTree)
			keys := pathT.Keys()
			fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
			if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
				for _, v := range keys {
					fmt.Printf("v: %s\n", v)
				}
				continue
			}
			for _, k := range keys {
				lPath := path + "." + k
				fmt.Printf("lPath: %s\n", lPath)
				lPathT := config.Get(lPath).(*toml.Tree) // TomlTree)
				fmt.Printf("\nlpath: %s tree keys: %#v\n", lPath, lPathT.Keys())
				for _, v := range lPathT.Keys() {
					fmt.Printf("v: %s\n", v)
				}
			}
		}
	}
	//TODO TESTME  fmt.Printf("Global configs all: %#v\n", gc)
	return &gc, config, nil
}

func checkNames(golden []string, toevaluate []string) error {
	if len(golden) != len(toevaluate) {
		e := fmt.Errorf("checkNames: names must be (%v) but are (%v)", golden, toevaluate)
		return e
	}
	for _, v := range toevaluate {
		if !keyIn(v, golden) {
			e := fmt.Errorf("checkNames: names must be (%v) but are (%v)", golden, toevaluate)
			return e
		}
	}
	return nil
}

//path := "iam"
func getKeys(path string, config *toml.Tree) []string {
	lkey := make([]string, 2)
	lkey[0] = "nokey"
	lkey[1] = ""
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
	matched, unknown := categoryMatchAndUnknowns(path, keys)
	if !matched || len(unknown) > 0 {
		if len(unknown) > 0 {
			fmt.Printf("unknown key list: %#v\n", unknown)
			keys = []string{"nokey"}
			for i := 0; i < len(unknown); i++ {
				keys = append(keys, unknown[i])
			}
			fmt.Printf("New key list: %#v\n", keys)
		}
		// These should be the unknown 'name' values
		//Need to find the order from within ie which one if Primary
		if path == "iam" {
			if len(keys) > 2 {
				//Can only have one or two
				fmt.Printf("TOO MANY IAM NAMES: %#v\n", keys)
			}

			prime := -1
			for i, v := range keys {
				lPath := path + "." + v + "." + "primary"
				if v == "nokey" {
					lPath = path + "." + "primary"
				}
				if config.Has(lPath) {
					lPathobj := config.Get(lPath)
					p := lPathobj.(bool)
					if p == true {
						if prime < 0 {
							prime = i
						} else {
							// ERRor Only one can be prime
						}
					}
				}
			}
			if prime < 0 {
				// Error At least one must be prime
			}
			lkey[0] = keys[prime]
			lkey[1] = keys[1]
			if prime == 1 {
				lkey[1] = keys[0]
			}
			return lkey
		}
		return keys
	}
	return lkey
}

//TODO FIXME THINKABOUTTHIS
// don't pass in the value from struct but rather use path to get value from toml Tree
func getPathStrValue(path string, config *toml.Tree) string {
	var targetVal interface{}
	obj := config.Get(path)
	switch reflect.TypeOf(obj).Name() {
	case "int64":
		targetVal = obj.(int64)
	case "float64":
		targetVal = obj.(float64)
	case "bool":
		targetVal = obj.(bool)
	case "string":
		targetVal = obj.(string)
	}
	s := fmt.Sprintf("%v", targetVal)
	return s
}

func setStringMapValueWithValue(ipsm *map[string]string,
	s string, val string) error {
	fmt.Printf("Set %s to %s\n", s, val)
	_, ok := (*ipsm)[s]
	if !ok {
		e := fmt.Errorf("setStringMapValue: attempt to set a non-existant parameter: %s", s)
		return e
	}
	(*ipsm)[s] = val
	return nil
}

func setStringMapValue(ipsm *map[string]string,
	s string, path string, config *toml.Tree) error {
	fmt.Printf("Attempting to set %s\n", s)
	_, ok := (*ipsm)[s]
	if !ok {
		e := fmt.Errorf("setStringMapValue: attempt to set a non-existant parameter: %s", s)
		return e
	}
	//(*ipsm)[s] = ""
	fmt.Printf("checking path: %s\n", path)
	Hval := -1
	indxH := strings.Index(path, "#")
	if indxH >= 0 {
		Hval = int(path[indxH+1] - '0')
		path = path[:indxH]
		fmt.Printf("now checking path: %s\n", path)
	}
	if config.Has(path) {
		v := getPathStrValue(path, config)
		fmt.Printf("DOES HAVE and it is: %s\n", v)
		if Hval > 0 {
			svals := strings.Split(v, "-")
			v = svals[Hval-1]
			fmt.Printf("DOES HAVE and will use: %s\n", v)
		}
		(*ipsm)[s] = v
	} else {
		fmt.Printf("DOES not have: %s\n", path)
	}
	return nil
}

func getInputStringsMapFromToml(filename string) map[string]string {
	/*gc*/ _, config, err := goGetTomlData(filename)
	if err != nil {
		fmt.Printf("getInputStringMapFromToml: %s\n", err)
		return nil
	}
	ipsm := rplanlib.NewInputStringsMap()

	//
	// Get all the unknown keys first
	//
	iamNames := getKeys("iam", config)
	fmt.Printf("iam names: %#v\n", iamNames)

	ssNames := getKeys("SocialSecurity", config)
	fmt.Printf("SS names: %#v\n", ssNames)
	err = checkNames(iamNames, ssNames)
	if err != nil {
		fmt.Printf("getInputStringMapFromToml: %s\n", err)
		return nil
	}
	iRANames := getKeys("IRA", config)
	fmt.Printf("IRA names: %#v\n", iRANames)
	err = checkNames(iamNames, iRANames)
	if err != nil {
		fmt.Printf("getInputStringMapFromToml: %s\n", err)
		return nil
	}
	rothNames := getKeys("roth", config)
	fmt.Printf("roth names: %#v\n", rothNames)
	err = checkNames(iamNames, rothNames)
	if err != nil {
		fmt.Printf("getInputStringMapFromToml: %s\n", err)
		return nil
	}
	aftertaxNames := getKeys("aftertax", config)
	fmt.Printf("aftertax names: %#v\n", aftertaxNames)
	assetNames := getKeys("asset", config)
	fmt.Printf("Asset names: %#v\n", assetNames)
	incomeNames := getKeys("income", config)
	fmt.Printf("Income names: %#v\n", incomeNames)
	expenseNames := getKeys("expense", config)
	fmt.Printf("Expense names: %#v\n", expenseNames)

	//
	// Now we can work our way though setting values in InputStrDefs
	//
	for i, k := range TomlStrDefs {
		if k[0] == '@' {
			// All keys used in TomlStrDefs are iam keys (names)
			indx := int(k[len(k)-1] - '0')
			n := iamNames[indx-1]
			err = setStringMapValueWithValue(&ipsm,
				rplanlib.InputStrDefs[i], n)
			if err != nil {
				fmt.Printf("getInputStringsMapFromToml: %s\n", err)
			}
			continue
		} else {
			indxP := strings.Index(k, "%")
			if indxP < 0 {
				err = setStringMapValue(&ipsm,
					rplanlib.InputStrDefs[i], k, config)
				if err != nil {
					fmt.Printf("getInputStringsMapFromToml: %s\n", err)
				}
				continue
			}
			val := k[indxP+1] - '0'
			fmt.Printf("Index val is %d\n", val)
			var p string
			fmt.Printf("*** iamNames[%d]: %s\n", val, iamNames[val])
			if iamNames[val] != "nokey" {
				p = strings.Replace(k, "%0", iamNames[val], 1)
				if val == 1 {
					p = strings.Replace(k, "%1", iamNames[val], 1)
				}
			} else {
				s := strings.Split(k, ".")
				p = s[0] + "." + s[2]
				fmt.Printf("have nokey using: %s\n", p)
			}
			fmt.Printf("Will use path: %s\n", p)

			err = setStringMapValue(&ipsm,
				rplanlib.InputStrDefs[i], p, config)
			if err != nil {
				fmt.Printf("getInputStringsMapFromToml: %s\n", err)
			}
			continue
		}
	}
	//
	// Now we can work our way though setting values in InputStreamStrDefs
	//
	var names []string
	for j := 1; j < rplanlib.MaxStreams+1; j++ {
		for i, k := range TomlStreamStrDefs {
			fmt.Printf("InputStrDefs[%d]: '%s', TomlStreamStrDefs[%d]: '%s'\n", i, rplanlib.InputStreamStrDefs[i], i, k)
			targetStr := fmt.Sprintf("%s%d", rplanlib.InputStreamStrDefs[i], j)
			if k[0] == '@' {
				switch k[1:] {
				case "income":
					names = incomeNames
				case "expense":
					names = expenseNames
				case "asset":
					names = assetNames
				default:
					fmt.Printf("EEEEEError 232323\n")
				}
				if len(names) < j {
					continue
				}
				n := names[j-1]
				err = setStringMapValueWithValue(&ipsm, targetStr, n)
				if err != nil {
					fmt.Printf("getInputStringsMapFromToml: %s\n", err)
				}
				continue
			} else { // TODO should not need else with the above continue
				indxP := strings.Index(k, "%")
				if indxP < 0 {
					err = setStringMapValue(&ipsm, targetStr, k, config)
					if err != nil {
						fmt.Printf("getInputStringsMapFromToml: %s\n", err)
					}
					continue
				}
				strs := strings.Split(k, ".")
				switch strs[0] {
				case "income":
					names = incomeNames
				case "expense":
					names = expenseNames
				case "asset":
					names = assetNames
				default:
					fmt.Printf("EEEEEError 8989\n")
				}
				if len(names) < j {
					continue
				}
				n := names[j-1]
				var p string
				fmt.Printf("*** names[%d]: %s\n", j-1, n)
				if n != "nokey" {
					p = strings.Replace(k, "%i", n, 1)
				} else {
					p = strs[0] + "." + strs[2]
					fmt.Printf("have nokey using: %s\n", p)
				}
				fmt.Printf("Will use path: %s\n", p)

				err = setStringMapValue(&ipsm, targetStr, p, config)
				if err != nil {
					fmt.Printf("getInputStringsMapFromToml: %s\n", err)
				}
				continue
			}
		}
	}
	return ipsm
}

/*
type GlobalConfig struct {
	Title          string
	RetirementType string `toml:"retirement_type"`
	Returns        float64
	Inflation      float64
	Maximize       string
	Iam            map[string]Iam
	SocialSecurity map[string]SocialSecurity
	IRA            map[string]IRA
	Roth           map[string]Roth
	Aftertax       Aftertax
	Asset          map[string]Asset
	Income         map[string]Income
	Expense        map[string]Expense
}
*/
