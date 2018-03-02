package main

//import "github.com/pelletier/go-toml"
import (
	"fmt"
	"reflect"

	"github.com/pelletier/go-toml"
)

type Iam struct {
	Primary                 bool
	Age                     int
	Retire                  int
	Through                 int
	DefinedContributionPlan string
}

var iamKeys = []string{"primary", "age", "retire", "through", "definedContributionPlan"}
var iamReqKeys = []string{"age", "retire", "through"}

type SocialSecurity struct {
	Amount int
	Age    string
}

var socialSecurityKeys = []string{"amount", "age"}
var socialSecurityReqKeys = socialSecurityKeys

type Income struct {
	Amount    int
	Age       string
	Inflation bool
	Tax       bool
}

var incomeKeys = []string{"amount", "age", "inflation", "tax"}
var incomeReqKeys = []string{"amount", "age"}

type Expense struct {
	Amount    int
	Age       string // FIXME TODO change toml def to be nicer eg this should be period rather than age
	Inflation bool
}

var expenseKeys = []string{"amount", "age", "inflation"}
var expenseReqKeys = []string{"amount", "age"}

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

type Desired struct {
	Amount int
}

var desiredKeys = []string{"amount"}
var desiredReqKeys = desiredKeys

type Max struct {
	Amount int
}

var maxKeys = []string{"amount"}
var maxReqKeys = maxKeys

type IRA struct {
	Bal       int
	Rate      float64
	Contrib   int
	Inflation bool
	Period    string
}

var IRAKeys = []string{"bal", "rate", "contrib", "inflation", "period"}
var IRAReqKeys = []string{"bal"}

type Roth struct {
	Bal       int
	Rate      float64
	Contrib   int
	Inflation bool
	Period    string
}

var rothKeys = IRAKeys
var rothReqKeys = IRAReqKeys

type Aftertax struct {
	Bal       int
	Basis     int
	Rate      float64
	Contrib   int
	Inflation bool
	Period    string
}

var aftertaxKeys = []string{"bal", "basis", "rate", "contrib", "inflation", "period"}
var aftertaxReqKeys = IRAReqKeys

type GlobalConfig struct {
	Title          string
	RetirementType string
	Returns        float64
	Inflation      float64
	Maximize       string
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

func doItWithUnMarshal() {
	//toml.Unmarshal(document, &person)
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
		fmt.Printf("aftertax config: %#v\n", aftertaxCfig)
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
			for _, v := range IRAKeys {
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
