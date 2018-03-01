package main

//import "github.com/pelletier/go-toml"
import (
	"fmt"
	"reflect"

	"github.com/pelletier/go-toml"
)

type iam struct {
	primary                 bool
	age                     int
	retire                  int
	through                 int
	definedContributionPlan string
}

var iamKeys = []string{"primary", "age", "retire", "through", "definedContributionPlan"}
var iamReqKeys = []string{"age", "retire", "through"}

type socialSecurity struct {
	amount int
	age    string
}

var socialSecurityKeys = []string{"amount", "age"}
var socialSecurityReqKeys = socialSecurityKeys

type income struct {
	amount    int
	age       string
	inflation bool
	tax       bool
}

var incomeKeys = []string{"amount", "age", "inflation", "tax"}
var incomeReqKeys = []string{"amount", "age"}

type expense struct {
	amount    int
	age       string // FIXME TODO change toml def to be nicer eg this should be period rather than age
	inflation bool
}

var expenseKeys = []string{"amount", "age", "inflation"}
var expenseReqKeys = []string{"amount", "age"}

type asset struct {
	value               int
	costAndImprovements int
	ageToSell           int
	owedAtAgeToSell     int
	primaryResidence    bool
	rate                float64
	brokerageRate       float64
}

var assetKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell", "primaryResidence", "rate", "brokerageRate"}
var assetReqKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell"}

type desired struct {
	amount int
}

var desiredKeys = []string{"amount"}
var desiredReqKeys = desiredKeys

type max struct {
	amount int
}

var maxKeys = []string{"amount"}
var maxReqKeys = maxKeys

type IRA struct {
	bal       int
	rate      float64
	contrib   int
	inflation bool
	period    string
}

var IRAKeys = []string{"bal", "rate", "contrib", "inflation", "period"}
var IRAReqKeys = []string{"bal"}

type roth struct {
	bal       int
	rate      float64
	contrib   int
	inflation bool
	period    string
}

var rothKeys = IRAKeys
var rothReqKeys = IRAReqKeys

type aftertax struct {
	bal       int
	basis     int
	rate      float64
	contrib   int
	inflation bool
	period    string
}

var aftertaxKeys = []string{"bal", "basis", "rate", "contrib", "inflation", "period"}
var aftertaxReqKeys = IRAReqKeys

type globalConfig struct {
	title          string
	retirementType string
	returns        float64
	inflation      float64
	maximize       string
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
	"income",
	"asset",
	"expense",
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
		reqKeys = IRAReqKeys
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

func getFloat64Value(obj interface{}) float64 {
	var targetVal float64
	switch reflect.TypeOf(obj).Name() {
	case "int64":
		targetVal = float64(obj.(int64))
	case "float64":
		targetVal = obj.(float64)
	}
	fmt.Printf("getFloat64Value, got: %g\n", targetVal)
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
	fmt.Printf("getInt64Value, got: %d\n", targetVal)
	return targetVal
}

func getTomlGlobal(config *toml.Tree) globalConfig {
	globals := globalConfig{}
	for _, lPath := range globalKeys {
		lPathobj := config.Get(lPath)
		switch lPath {
		case "title":
			if config.Has(lPath) {
				globals.title = lPathobj.(string)
				continue
			}
			globals.title = ""
		case "retirement_type":
			if config.Has(lPath) {
				globals.retirementType = lPathobj.(string)
				continue
			}
			globals.retirementType = ""
		case "returns":
			if config.Has(lPath) {
				globals.returns = getFloat64Value(lPathobj)
				continue
			}
			globals.returns = -1
		case "inflation":
			if config.Has(lPath) {
				globals.inflation = getFloat64Value(lPathobj)
				continue
			}
			globals.inflation = -1
		case "maximize":
			if config.Has(lPath) {
				globals.maximize = lPathobj.(string)
				continue
			}
			globals.maximize = ""
		default:
			fmt.Printf("*** Implement ME(%s) ****\n", lPath)
		}
	}
	return globals
}

func getTomlAftertax(config *toml.Tree, path string) *aftertax {

	if config.Has(path) {
		aftertaxCfig := aftertax{}
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
						aftertaxCfig.bal = int(getInt64Value(lPathobj))
						continue
					}
					aftertaxCfig.bal = -1
				case "basis":
					if config.Has(lPath) {
						aftertaxCfig.basis = int(getInt64Value(lPathobj))
						continue
					}
					aftertaxCfig.basis = -1
				case "contrib":
					if config.Has(lPath) {
						aftertaxCfig.contrib = int(getInt64Value(lPathobj))
						continue
					}
					aftertaxCfig.contrib = -1
				case "rate":
					if config.Has(lPath) {
						aftertaxCfig.rate = getFloat64Value(lPathobj)
						continue
					}
					aftertaxCfig.rate = -1
				case "inflation":
					if config.Has(lPath) {
						aftertaxCfig.inflation = lPathobj.(bool)
						continue
					}
					aftertaxCfig.inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "period":
					if config.Has(lPath) {
						aftertaxCfig.period = lPathobj.(string)
						continue
					}
					aftertaxCfig.period = ""
				}
			}
		} else {
			fmt.Printf("*** EERROOOORR ***\n")
			return nil
		}
		fmt.Printf("aftertax config: %#v\n", aftertaxCfig)
		return &aftertaxCfig
	}
	return nil
}

func getTomlIAm(config *toml.Tree, path string) *iam {

	if config.Has(path) {
		iamCfig := iam{}
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
						iamCfig.age = int(getInt64Value(lPathobj))
						continue
					}
					iamCfig.age = -1
				case "retire":
					if config.Has(lPath) {
						iamCfig.retire = int(getInt64Value(lPathobj))
						continue
					}
					iamCfig.retire = -1
				case "through":
					if config.Has(lPath) {
						iamCfig.through = int(getInt64Value(lPathobj))
						continue
					}
					iamCfig.through = -1
				case "definedContributionPlan":
					if config.Has(lPath) {
						iamCfig.definedContributionPlan = lPathobj.(string)
						continue
					}
					iamCfig.definedContributionPlan = ""
				case "primary":
					if config.Has(lPath) {
						iamCfig.primary = lPathobj.(bool)
						continue
					}
					iamCfig.primary = false
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

func getTomlIAmMap(config *toml.Tree, path string) *map[string]iam {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	//fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
	m := make(map[string]iam)
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

func getTomlSocialSecurity(config *toml.Tree, path string) *socialSecurity {

	if config.Has(path) {
		ssCfig := socialSecurity{}
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
						ssCfig.amount = int(getInt64Value(lPathobj))
						continue
					}
					ssCfig.amount = -1
				case "age":
					if config.Has(lPath) {
						ssCfig.age = lPathobj.(string)
						continue
					}
					ssCfig.age = ""
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

func getTomlSocialSecurityMap(config *toml.Tree, path string) *map[string]socialSecurity {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	//fmt.Printf("\npath: %s tree keys: %#v\n", path, keys)
	m := make(map[string]socialSecurity)
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
			for _, v := range IRAKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				//var IRAKeys = []string{"bal", "rate", "contrib", "inflation", "period"}
				switch v {
				case "bal":
					if config.Has(lPath) {
						iraCfig.bal = int(getInt64Value(lPathobj))
						continue
					}
					iraCfig.bal = -1
				case "rate":
					if config.Has(lPath) {
						iraCfig.rate = getFloat64Value(lPathobj)
						continue
					}
					iraCfig.rate = -1
				case "contrib":
					if config.Has(lPath) {
						iraCfig.contrib = int(getInt64Value(lPathobj))
						continue
					}
					iraCfig.contrib = -1
				case "inflation":
					if config.Has(lPath) {
						iraCfig.inflation = lPathobj.(bool)
						continue
					}
					iraCfig.inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "period":
					if config.Has(lPath) {
						iraCfig.period = lPathobj.(string)
						continue
					}
					iraCfig.period = ""
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

func getTomlRoth(config *toml.Tree, path string) *roth {

	if config.Has(path) {
		rothCfig := roth{}
		pathT := config.Get(path).(*toml.Tree)
		keys := pathT.Keys()
		if categoryMatch(path, keys) { // FIXME TODO maybe this should be a leaf check instead
			for _, v := range rothKeys {
				lPath := path + "." + v
				lPathobj := config.Get(lPath)
				switch v {
				case "bal":
					if config.Has(lPath) {
						rothCfig.bal = int(getInt64Value(lPathobj))
						continue
					}
					rothCfig.bal = -1
				case "rate":
					if config.Has(lPath) {
						rothCfig.rate = getFloat64Value(lPathobj)
						continue
					}
					rothCfig.rate = -1
				case "contrib":
					if config.Has(lPath) {
						rothCfig.contrib = int(getInt64Value(lPathobj))
						continue
					}
					rothCfig.contrib = -1
				case "inflation":
					if config.Has(lPath) {
						rothCfig.inflation = lPathobj.(bool)
						continue
					}
					rothCfig.inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "period":
					if config.Has(lPath) {
						rothCfig.period = lPathobj.(string)
						continue
					}
					rothCfig.period = ""
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

func getTomlRothMap(config *toml.Tree, path string) *map[string]roth {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]roth)
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

func getTomlIncome(config *toml.Tree, path string) *income {

	if config.Has(path) {
		iCfig := income{}
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
						iCfig.amount = int(getInt64Value(lPathobj))
						continue
					}
					iCfig.amount = -1
				case "inflation":
					if config.Has(lPath) {
						iCfig.inflation = lPathobj.(bool)
						continue
					}
					iCfig.inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "tax":
					if config.Has(lPath) {
						iCfig.tax = lPathobj.(bool)
						continue
					}
					iCfig.tax = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "age":
					if config.Has(lPath) {
						iCfig.age = lPathobj.(string)
						continue
					}
					iCfig.age = ""
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

func getTomlIncomeMap(config *toml.Tree, path string) *map[string]income {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]income)
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

func getTomlExpense(config *toml.Tree, path string) *expense {

	if config.Has(path) {
		eCfig := expense{}
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
						eCfig.amount = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.amount = -1
				case "inflation":
					if config.Has(lPath) {
						eCfig.inflation = lPathobj.(bool)
						continue
					}
					eCfig.inflation = true
					//TODO FIXME setting default bool to TRUE IS THIS WHAT IT SHOULD BE???
				case "age":
					if config.Has(lPath) {
						eCfig.age = lPathobj.(string)
						continue
					}
					eCfig.age = ""
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

func getTomlExpenseMap(config *toml.Tree, path string) *map[string]expense {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]expense)
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

func getTomlAsset(config *toml.Tree, path string) *asset {

	if config.Has(path) {
		eCfig := asset{}
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
						eCfig.value = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.value = -1
				case "primaryResidence":
					if config.Has(lPath) {
						eCfig.primaryResidence = lPathobj.(bool)
						continue
					}
					eCfig.primaryResidence = false
					//TODO FIXME setting default bool to false IS THIS WHAT IT SHOULD BE???
				case "costAndImprovements":
					if config.Has(lPath) {
						eCfig.costAndImprovements = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.costAndImprovements = -1
				case "ageToSell":
					if config.Has(lPath) {
						eCfig.ageToSell = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.ageToSell = -1
				case "owedAtAgeToSell":
					if config.Has(lPath) {
						eCfig.owedAtAgeToSell = int(getInt64Value(lPathobj))
						continue
					}
					eCfig.owedAtAgeToSell = -1
				case "rate":
					if config.Has(lPath) {
						eCfig.rate = getFloat64Value(lPathobj)
						continue
					}
					eCfig.rate = -1
				//var assetKeys = []string{"value", "costAndImprovements", "ageToSell", "owedAtAgeToSell", "primaryResidence", "rate", "broderageRate"}
				case "brokerageRate":
					if config.Has(lPath) {
						eCfig.brokerageRate = getFloat64Value(lPathobj)
						continue
					}
					eCfig.brokerageRate = -1
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

func getTomlAssetMap(config *toml.Tree, path string) *map[string]asset {
	pathT := config.Get(path).(*toml.Tree)
	keys := pathT.Keys()
	m := make(map[string]asset)
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

func goGetTomlData() {

	config, err := toml.LoadFile("mobile_j.toml")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	//fmt.Printf("\nKeys: %#v\n", config.Keys())
	//for _, b := range config.Keys()
	for _, path := range categoryKeys {
		switch path {
		case "global":
			c := getTomlGlobal(config)
			fmt.Printf("global configs: %#v\n", c)
		case "iam":
			c := getTomlIAmMap(config, path)
			fmt.Printf("iam configs map: %#v\n", c)
		case "SocialSecurity":
			c := getTomlSocialSecurityMap(config, path)
			fmt.Printf("socail security configs map: %#v\n", c)
		case "IRA":
			c := getTomlIRAMap(config, path)
			fmt.Printf("IRA configs map: %#v\n", c)
		case "roth":
			c := getTomlRothMap(config, path)
			fmt.Printf("Roth configs map: %#v\n", c)
		case "aftertax":
			c := getTomlAftertax(config, path)
			fmt.Printf("aftertax configs: %#v\n", c)
		case "income":
			c := getTomlIncomeMap(config, path)
			fmt.Printf("Income configs map: %#v\n", c)
		case "asset":
			c := getTomlAssetMap(config, path)
			fmt.Printf("Asset configs map: %#v\n", c)
		case "expense":
			c := getTomlExpenseMap(config, path)
			fmt.Printf("Expense configs map: %#v\n", c)
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
}

func try2() {
	user := "not me"
	password := "password"

	config, _ := toml.Load(`
		[postgres]
		user = "pelletier"
		password = "mypassword"`)
	// retrieve data directly
	user = config.Get("postgres.user").(string)

	// or using an intermediate object
	postgresConfig := config.Get("postgres").(*toml.Tree) //TomlTree)
	password = postgresConfig.Get("password").(string)

	fmt.Printf("\nTry2:\n")
	fmt.Printf("\nuser: %s, password: %s\n", user, password)
}
