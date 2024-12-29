// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import (
	"strconv"
	"strings"
)

// ISO3166CountryCode associates ISO 3166-1 alpha-2, alpha-3, and numeric codes
// with DXCC entities.  Most country codes are associated with a single DXCC
// entity via the Country array.  Some ISO 3166 countries are made up of
// several smaller DXCC entities like the Russian Federation (European Russia
// and Asiatic Russia plus Kaliningrad) and the United Kingdom (England, Wales,
// Scotland, and Northern Ireland).  Some DXCC entities are primary subdivisions
// within a country-level DXCC entity, such as Alaska and Hawaii with the
// United States of America, Sardinia with Italy, or Andaman & Nicobar Islands
// with India; in such cases the parent is the first entity in the DXCC array
// and the others follow.  Some DXCC entities are not associated with an ISO
// country code; these are mostly remote islands without permanent civilian
// populations, along with special entities like ITU Headquarters.  Deleted DXCC
// entities like Czechoslovakia are not associated with a country code.
// Likewise, deleted ISO 3166-1 codes are not present.  Geopolitical events and
// administrative updates by ISO or ARRL may lead to a change in this data;
// such changes will not be considered a semantic versioning breaking change.
type ISO3166CountryCode struct {
	Alpha2       string                 // ISO 3166-1 alpha-2 code (two ASCII letters)
	Alpha3       string                 // ISO 3166-1 alpha-3 code (three ASCII letters)
	Numeric      string                 // ISO 3166-1 numeric code (three ASCII digits)
	EnglishName  string                 // Official English name of the country, mixed case
	DXCC         []CountryEnum          // One or more DXCC entities that make up this country
	Subdivisions map[string]CountryEnum // Subdivision codes associated with specific sub-national DXCC entities
}

func (c ISO3166CountryCode) IncludesDXCC(dxcc string) bool {
	for _, d := range c.DXCC {
		if d.EntityCode == dxcc || strings.EqualFold(d.EntityName, dxcc) {
			return true
		}
	}
	return false
}

var (
	// ISO3166Alpha maps ISO 3166-1 alpha-2 and alpha-3 codes to countries.
	ISO3166Alpha = make(map[string]ISO3166CountryCode)
	// ISO3166Alpha maps ISO 3166-1 numeric codes to countries.
	ISO3166Numeric = make(map[int]ISO3166CountryCode)
)

func init() {
	for _, c := range ISO3166Countries {
		ISO3166Alpha[c.Alpha2] = c
		ISO3166Alpha[c.Alpha3] = c
		if n, err := strconv.Atoi(c.Numeric); err == nil {
			ISO3166Numeric[n] = c
		}
	}
}

var (
	// TODO Change Kosovo details when it gets an official assigned ISO code.
	CountryCodeXKX = ISO3166CountryCode{
		Alpha2:      "XK",  // see https://en.wikipedia.org/wiki/XK_(user_assigned_code)
		Alpha3:      "XKX", // some people have used XXK too
		Numeric:     "999",
		EnglishName: "Kosovo (Republic of)",
		DXCC:        []CountryEnum{CountryRepublicOfKosovo},
	}

	CountryCodeAFG = ISO3166CountryCode{
		Alpha2:      "AF",
		Alpha3:      "AFG",
		Numeric:     "004",
		EnglishName: "Afghanistan",
		DXCC:        []CountryEnum{CountryAfghanistan},
	}

	CountryCodeALB = ISO3166CountryCode{
		Alpha2:      "AL",
		Alpha3:      "ALB",
		Numeric:     "008",
		EnglishName: "Albania",
		DXCC:        []CountryEnum{CountryAlbania},
	}

	CountryCodeDZA = ISO3166CountryCode{
		Alpha2:      "DZ",
		Alpha3:      "DZA",
		Numeric:     "012",
		EnglishName: "Algeria",
		DXCC:        []CountryEnum{CountryAlgeria},
	}

	CountryCodeASM = ISO3166CountryCode{
		Alpha2:      "AS",
		Alpha3:      "ASM",
		Numeric:     "016",
		EnglishName: "American Samoa",
		DXCC:        []CountryEnum{CountryAmericanSamoa},
	}

	CountryCodeAND = ISO3166CountryCode{
		Alpha2:      "AD",
		Alpha3:      "AND",
		Numeric:     "020",
		EnglishName: "Andorra",
		DXCC:        []CountryEnum{CountryAndorra},
	}

	CountryCodeAGO = ISO3166CountryCode{
		Alpha2:      "AO",
		Alpha3:      "AGO",
		Numeric:     "024",
		EnglishName: "Angola",
		DXCC:        []CountryEnum{CountryAngola},
	}

	CountryCodeAIA = ISO3166CountryCode{
		Alpha2:      "AI",
		Alpha3:      "AIA",
		Numeric:     "660",
		EnglishName: "Anguilla",
		DXCC:        []CountryEnum{CountryAnguilla},
	}

	CountryCodeATA = ISO3166CountryCode{
		Alpha2:      "AQ",
		Alpha3:      "ATA",
		Numeric:     "010",
		EnglishName: "Antarctica",
		DXCC:        []CountryEnum{CountryAntarctica},
	}

	CountryCodeATG = ISO3166CountryCode{
		Alpha2:      "AG",
		Alpha3:      "ATG",
		Numeric:     "028",
		EnglishName: "Antigua and Barbuda",
		DXCC:        []CountryEnum{CountryAntiguaBarbuda},
	}

	CountryCodeARG = ISO3166CountryCode{
		Alpha2:      "AR",
		Alpha3:      "ARG",
		Numeric:     "032",
		EnglishName: "Argentina",
		DXCC:        []CountryEnum{CountryArgentina},
	}

	CountryCodeARM = ISO3166CountryCode{
		Alpha2:      "AM",
		Alpha3:      "ARM",
		Numeric:     "051",
		EnglishName: "Armenia",
		DXCC:        []CountryEnum{CountryArmenia},
	}

	CountryCodeABW = ISO3166CountryCode{
		Alpha2:      "AW",
		Alpha3:      "ABW",
		Numeric:     "533",
		EnglishName: "Aruba",
		DXCC:        []CountryEnum{CountryAruba},
	}

	CountryCodeAUS = ISO3166CountryCode{
		Alpha2:      "AU",
		Alpha3:      "AUS",
		Numeric:     "036",
		EnglishName: "Australia",
		DXCC:        []CountryEnum{CountryAustralia},
	}

	CountryCodeAUT = ISO3166CountryCode{
		Alpha2:      "AT",
		Alpha3:      "AUT",
		Numeric:     "040",
		EnglishName: "Austria",
		DXCC:        []CountryEnum{CountryAustria},
	}

	CountryCodeAZE = ISO3166CountryCode{
		Alpha2:      "AZ",
		Alpha3:      "AZE",
		Numeric:     "031",
		EnglishName: "Azerbaijan",
		DXCC:        []CountryEnum{CountryAzerbaijan},
	}

	CountryCodeBHS = ISO3166CountryCode{
		Alpha2:      "BS",
		Alpha3:      "BHS",
		Numeric:     "044",
		EnglishName: "Bahamas (the)",
		DXCC:        []CountryEnum{CountryBahamas},
	}

	CountryCodeBHR = ISO3166CountryCode{
		Alpha2:      "BH",
		Alpha3:      "BHR",
		Numeric:     "048",
		EnglishName: "Bahrain",
		DXCC:        []CountryEnum{CountryBahrain},
	}

	CountryCodeBGD = ISO3166CountryCode{
		Alpha2:      "BD",
		Alpha3:      "BGD",
		Numeric:     "050",
		EnglishName: "Bangladesh",
		DXCC:        []CountryEnum{CountryBangladesh},
	}

	CountryCodeBRB = ISO3166CountryCode{
		Alpha2:      "BB",
		Alpha3:      "BRB",
		Numeric:     "052",
		EnglishName: "Barbados",
		DXCC:        []CountryEnum{CountryBarbados},
	}

	CountryCodeBLR = ISO3166CountryCode{
		Alpha2:      "BY",
		Alpha3:      "BLR",
		Numeric:     "112",
		EnglishName: "Belarus",
		DXCC:        []CountryEnum{CountryBelarus},
	}

	CountryCodeBEL = ISO3166CountryCode{
		Alpha2:      "BE",
		Alpha3:      "BEL",
		Numeric:     "056",
		EnglishName: "Belgium",
		DXCC:        []CountryEnum{CountryBelgium},
	}

	CountryCodeBLZ = ISO3166CountryCode{
		Alpha2:      "BZ",
		Alpha3:      "BLZ",
		Numeric:     "084",
		EnglishName: "Belize",
		DXCC:        []CountryEnum{CountryBelize},
	}

	CountryCodeBEN = ISO3166CountryCode{
		Alpha2:      "BJ",
		Alpha3:      "BEN",
		Numeric:     "204",
		EnglishName: "Benin",
		DXCC:        []CountryEnum{CountryBenin},
	}

	CountryCodeBMU = ISO3166CountryCode{
		Alpha2:      "BM",
		Alpha3:      "BMU",
		Numeric:     "060",
		EnglishName: "Bermuda",
		DXCC:        []CountryEnum{CountryBermuda},
	}

	CountryCodeBTN = ISO3166CountryCode{
		Alpha2:      "BT",
		Alpha3:      "BTN",
		Numeric:     "064",
		EnglishName: "Bhutan",
		DXCC:        []CountryEnum{CountryBhutan},
	}

	CountryCodeBOL = ISO3166CountryCode{
		Alpha2:      "BO",
		Alpha3:      "BOL",
		Numeric:     "068",
		EnglishName: "Bolivia (Plurinational State of)",
		DXCC:        []CountryEnum{CountryBolivia},
	}

	CountryCodeBES = ISO3166CountryCode{
		Alpha2:      "BQ",
		Alpha3:      "BES",
		Numeric:     "535",
		EnglishName: "Bonaire, Sint Eustatius and Saba",
		DXCC:        []CountryEnum{CountryBonaire, CountrySabaStEustatius},
		Subdivisions: map[string]CountryEnum{
			"BO": CountryBonaire,
			"SA": CountrySabaStEustatius,
			"SE": CountrySabaStEustatius,
		},
	}

	CountryCodeBIH = ISO3166CountryCode{
		Alpha2:      "BA",
		Alpha3:      "BIH",
		Numeric:     "070",
		EnglishName: "Bosnia and Herzegovina",
		DXCC:        []CountryEnum{CountryBosniaHerzegovina},
	}

	CountryCodeBWA = ISO3166CountryCode{
		Alpha2:      "BW",
		Alpha3:      "BWA",
		Numeric:     "072",
		EnglishName: "Botswana",
		DXCC:        []CountryEnum{CountryBotswana},
	}

	CountryCodeBVT = ISO3166CountryCode{
		Alpha2:      "BV",
		Alpha3:      "BVT",
		Numeric:     "074",
		EnglishName: "Bouvet Island",
		DXCC:        []CountryEnum{CountryBouvet},
	}

	CountryCodeBRA = ISO3166CountryCode{
		Alpha2:      "BR",
		Alpha3:      "BRA",
		Numeric:     "076",
		EnglishName: "Brazil",
		DXCC:        []CountryEnum{CountryBrazil},
	}

	CountryCodeIOT = ISO3166CountryCode{
		Alpha2:      "IO",
		Alpha3:      "IOT",
		Numeric:     "086",
		EnglishName: "British Indian Ocean Territory (the)",
		DXCC:        []CountryEnum{CountryChagosIslands},
	}

	CountryCodeBRN = ISO3166CountryCode{
		Alpha2:      "BN",
		Alpha3:      "BRN",
		Numeric:     "096",
		EnglishName: "Brunei Darussalam",
		DXCC:        []CountryEnum{CountryBruneiDarussalam},
	}

	CountryCodeBGR = ISO3166CountryCode{
		Alpha2:      "BG",
		Alpha3:      "BGR",
		Numeric:     "100",
		EnglishName: "Bulgaria",
		DXCC:        []CountryEnum{CountryBulgaria},
	}

	CountryCodeBFA = ISO3166CountryCode{
		Alpha2:      "BF",
		Alpha3:      "BFA",
		Numeric:     "854",
		EnglishName: "Burkina Faso",
		DXCC:        []CountryEnum{CountryBurkinaFaso},
	}

	CountryCodeBDI = ISO3166CountryCode{
		Alpha2:      "BI",
		Alpha3:      "BDI",
		Numeric:     "108",
		EnglishName: "Burundi",
		DXCC:        []CountryEnum{CountryBurundi},
	}

	CountryCodeCPV = ISO3166CountryCode{
		Alpha2:      "CV",
		Alpha3:      "CPV",
		Numeric:     "132",
		EnglishName: "Cabo Verde",
		DXCC:        []CountryEnum{CountryCapeVerde},
	}

	CountryCodeKHM = ISO3166CountryCode{
		Alpha2:      "KH",
		Alpha3:      "KHM",
		Numeric:     "116",
		EnglishName: "Cambodia",
		DXCC:        []CountryEnum{CountryCambodia},
	}

	CountryCodeCMR = ISO3166CountryCode{
		Alpha2:      "CM",
		Alpha3:      "CMR",
		Numeric:     "120",
		EnglishName: "Cameroon",
		DXCC:        []CountryEnum{CountryCameroon},
	}

	CountryCodeCAN = ISO3166CountryCode{
		Alpha2:      "CA",
		Alpha3:      "CAN",
		Numeric:     "124",
		EnglishName: "Canada",
		DXCC:        []CountryEnum{CountryCanada},
	}

	CountryCodeCYM = ISO3166CountryCode{
		Alpha2:      "KY",
		Alpha3:      "CYM",
		Numeric:     "136",
		EnglishName: "Cayman Islands (the)",
		DXCC:        []CountryEnum{CountryCaymanIslands},
	}

	CountryCodeCAF = ISO3166CountryCode{
		Alpha2:      "CF",
		Alpha3:      "CAF",
		Numeric:     "140",
		EnglishName: "Central African Republic (the)",
		DXCC:        []CountryEnum{CountryCentralAfrica},
	}

	CountryCodeTCD = ISO3166CountryCode{
		Alpha2:      "TD",
		Alpha3:      "TCD",
		Numeric:     "148",
		EnglishName: "Chad",
		DXCC:        []CountryEnum{CountryChad},
	}

	CountryCodeCHL = ISO3166CountryCode{
		Alpha2:      "CL",
		Alpha3:      "CHL",
		Numeric:     "152",
		EnglishName: "Chile",
		DXCC:        []CountryEnum{CountryChile},
	}

	CountryCodeCHN = ISO3166CountryCode{
		Alpha2:      "CN",
		Alpha3:      "CHN",
		Numeric:     "156",
		EnglishName: "China",
		DXCC:        []CountryEnum{CountryChina},
	}

	CountryCodeCXR = ISO3166CountryCode{
		Alpha2:      "CX",
		Alpha3:      "CXR",
		Numeric:     "162",
		EnglishName: "Christmas Island",
		DXCC:        []CountryEnum{CountryChristmasIsland},
	}

	CountryCodeCCK = ISO3166CountryCode{
		Alpha2:      "CC",
		Alpha3:      "CCK",
		Numeric:     "166",
		EnglishName: "Cocos (Keeling) Islands (the)",
		DXCC:        []CountryEnum{CountryCocosKeelingIslands},
	}

	CountryCodeCOL = ISO3166CountryCode{
		Alpha2:       "CO",
		Alpha3:       "COL",
		Numeric:      "170",
		EnglishName:  "Colombia",
		DXCC:         []CountryEnum{CountryColombia, CountrySanAndresProvidencia},
		Subdivisions: map[string]CountryEnum{"SAP": CountrySanAndresProvidencia},
	}

	CountryCodeCOM = ISO3166CountryCode{
		Alpha2:      "KM",
		Alpha3:      "COM",
		Numeric:     "174",
		EnglishName: "Comoros (the)",
		DXCC:        []CountryEnum{CountryComoros},
	}

	CountryCodeCOD = ISO3166CountryCode{
		Alpha2:      "CD",
		Alpha3:      "COD",
		Numeric:     "180",
		EnglishName: "Congo (the Democratic Republic of the)",
		DXCC:        []CountryEnum{CountryDemocraticRepublicOfTheCongo},
	}

	CountryCodeCOG = ISO3166CountryCode{
		Alpha2:      "CG",
		Alpha3:      "COG",
		Numeric:     "178",
		EnglishName: "Congo (the)",
		DXCC:        []CountryEnum{CountryRepublicOfTheCongo},
	}

	CountryCodeCOK = ISO3166CountryCode{
		Alpha2:      "CK",
		Alpha3:      "COK",
		Numeric:     "184",
		EnglishName: "Cook Islands (the)",
		DXCC:        []CountryEnum{CountrySouthCookIslands, CountryNorthCookIslands},
	}

	CountryCodeCRI = ISO3166CountryCode{
		Alpha2:      "CR",
		Alpha3:      "CRI",
		Numeric:     "188",
		EnglishName: "Costa Rica",
		DXCC:        []CountryEnum{CountryCostaRica},
	}

	CountryCodeHRV = ISO3166CountryCode{
		Alpha2:      "HR",
		Alpha3:      "HRV",
		Numeric:     "191",
		EnglishName: "Croatia",
		DXCC:        []CountryEnum{CountryCroatia},
	}

	CountryCodeCUB = ISO3166CountryCode{
		Alpha2:      "CU",
		Alpha3:      "CUB",
		Numeric:     "192",
		EnglishName: "Cuba",
		DXCC:        []CountryEnum{CountryCuba},
	}

	CountryCodeCUW = ISO3166CountryCode{
		Alpha2:      "CW",
		Alpha3:      "CUW",
		Numeric:     "531",
		EnglishName: "Curaçao",
		DXCC:        []CountryEnum{CountryCuracao},
	}

	CountryCodeCYP = ISO3166CountryCode{
		Alpha2:      "CY",
		Alpha3:      "CYP",
		Numeric:     "196",
		EnglishName: "Cyprus",
		DXCC:        []CountryEnum{CountryCyprus},
	}

	CountryCodeCZE = ISO3166CountryCode{
		Alpha2:      "CZ",
		Alpha3:      "CZE",
		Numeric:     "203",
		EnglishName: "Czechia",
		DXCC:        []CountryEnum{CountryCzechRepublic},
	}

	CountryCodeCIV = ISO3166CountryCode{
		Alpha2:      "CI",
		Alpha3:      "CIV",
		Numeric:     "384",
		EnglishName: "Côte d'Ivoire",
		DXCC:        []CountryEnum{CountryCoteDIvoire},
	}

	CountryCodeDNK = ISO3166CountryCode{
		Alpha2:      "DK",
		Alpha3:      "DNK",
		Numeric:     "208",
		EnglishName: "Denmark",
		DXCC:        []CountryEnum{CountryDenmark},
	}

	CountryCodeDJI = ISO3166CountryCode{
		Alpha2:      "DJ",
		Alpha3:      "DJI",
		Numeric:     "262",
		EnglishName: "Djibouti",
		DXCC:        []CountryEnum{CountryDjibouti},
	}

	CountryCodeDMA = ISO3166CountryCode{
		Alpha2:      "DM",
		Alpha3:      "DMA",
		Numeric:     "212",
		EnglishName: "Dominica",
		DXCC:        []CountryEnum{CountryDominica},
	}

	CountryCodeDOM = ISO3166CountryCode{
		Alpha2:      "DO",
		Alpha3:      "DOM",
		Numeric:     "214",
		EnglishName: "Dominican Republic (the)",
		DXCC:        []CountryEnum{CountryDominicanRepublic},
	}

	CountryCodeECU = ISO3166CountryCode{
		Alpha2:       "EC",
		Alpha3:       "ECU",
		Numeric:      "218",
		EnglishName:  "Ecuador",
		DXCC:         []CountryEnum{CountryEcuador, CountryGalapagosIslands},
		Subdivisions: map[string]CountryEnum{"W": CountryGalapagosIslands},
	}

	CountryCodeEGY = ISO3166CountryCode{
		Alpha2:      "EG",
		Alpha3:      "EGY",
		Numeric:     "818",
		EnglishName: "Egypt",
		DXCC:        []CountryEnum{CountryEgypt},
	}

	CountryCodeSLV = ISO3166CountryCode{
		Alpha2:      "SV",
		Alpha3:      "SLV",
		Numeric:     "222",
		EnglishName: "El Salvador",
		DXCC:        []CountryEnum{CountryElSalvador},
	}

	CountryCodeGNQ = ISO3166CountryCode{
		Alpha2:       "GQ",
		Alpha3:       "GNQ",
		Numeric:      "226",
		EnglishName:  "Equatorial Guinea",
		DXCC:         []CountryEnum{CountryEquatorialGuinea, CountryAnnobonIsland},
		Subdivisions: map[string]CountryEnum{"AN": CountryAnnobonIsland},
	}

	CountryCodeERI = ISO3166CountryCode{
		Alpha2:      "ER",
		Alpha3:      "ERI",
		Numeric:     "232",
		EnglishName: "Eritrea",
		DXCC:        []CountryEnum{CountryEritrea},
	}

	CountryCodeEST = ISO3166CountryCode{
		Alpha2:      "EE",
		Alpha3:      "EST",
		Numeric:     "233",
		EnglishName: "Estonia",
		DXCC:        []CountryEnum{CountryEstonia},
	}

	CountryCodeSWZ = ISO3166CountryCode{
		Alpha2:      "SZ",
		Alpha3:      "SWZ",
		Numeric:     "748",
		EnglishName: "Eswatini",
		DXCC:        []CountryEnum{CountryKingdomOfEswatini},
	}

	CountryCodeETH = ISO3166CountryCode{
		Alpha2:      "ET",
		Alpha3:      "ETH",
		Numeric:     "231",
		EnglishName: "Ethiopia",
		DXCC:        []CountryEnum{CountryEthiopia},
	}

	CountryCodeFLK = ISO3166CountryCode{
		Alpha2:      "FK",
		Alpha3:      "FLK",
		Numeric:     "238",
		EnglishName: "Falkland Islands (the) [Malvinas]",
		DXCC:        []CountryEnum{CountryFalklandIslands},
	}

	CountryCodeFRO = ISO3166CountryCode{
		Alpha2:      "FO",
		Alpha3:      "FRO",
		Numeric:     "234",
		EnglishName: "Faroe Islands (the)",
		DXCC:        []CountryEnum{CountryFaroeIslands},
	}

	CountryCodeFJI = ISO3166CountryCode{
		Alpha2:       "FJ",
		Alpha3:       "FJI",
		Numeric:      "242",
		EnglishName:  "Fiji",
		DXCC:         []CountryEnum{CountryFiji, CountryRotumaIsland},
		Subdivisions: map[string]CountryEnum{"R": CountryRotumaIsland},
	}

	CountryCodeFIN = ISO3166CountryCode{
		Alpha2:      "FI",
		Alpha3:      "FIN",
		Numeric:     "246",
		EnglishName: "Finland",
		DXCC:        []CountryEnum{CountryFinland},
	}

	CountryCodeFRA = ISO3166CountryCode{
		Alpha2:      "FR",
		Alpha3:      "FRA",
		Numeric:     "250",
		EnglishName: "France",
		DXCC:        []CountryEnum{CountryFrance, CountryCorsica},
		Subdivisions: map[string]CountryEnum{
			"2A":  CountryCorsica, // Corse-du-Sud
			"2B":  CountryCorsica, // Haute-Corse
			"20R": CountryCorsica, // Corsica as a metropolitan collectivity
		},
	}

	CountryCodeGUF = ISO3166CountryCode{
		Alpha2:      "GF",
		Alpha3:      "GUF",
		Numeric:     "254",
		EnglishName: "French Guiana",
		DXCC:        []CountryEnum{CountryFrenchGuiana},
	}

	CountryCodePYF = ISO3166CountryCode{
		Alpha2:      "PF",
		Alpha3:      "PYF",
		Numeric:     "258",
		EnglishName: "French Polynesia",
		DXCC:        []CountryEnum{CountryFrenchPolynesia, CountryMarquesasIslands, CountryAustralIsland},
	}

	CountryCodeATF = ISO3166CountryCode{
		Alpha2:      "TF",
		Alpha3:      "ATF",
		Numeric:     "260",
		EnglishName: "French Southern Territories (the)",
		DXCC:        []CountryEnum{CountryKerguelenIslands, CountryCrozetIsland, CountryAmsterdamStPaulIslands, CountryJuanDeNovaEuropa, CountryTromelinIsland, CountryGloriosoIslands},
	}

	CountryCodeGAB = ISO3166CountryCode{
		Alpha2:      "GA",
		Alpha3:      "GAB",
		Numeric:     "266",
		EnglishName: "Gabon",
		DXCC:        []CountryEnum{CountryGabon},
	}

	CountryCodeGMB = ISO3166CountryCode{
		Alpha2:      "GM",
		Alpha3:      "GMB",
		Numeric:     "270",
		EnglishName: "Gambia (the)",
		DXCC:        []CountryEnum{CountryTheGambia},
	}

	CountryCodeGEO = ISO3166CountryCode{
		Alpha2:      "GE",
		Alpha3:      "GEO",
		Numeric:     "268",
		EnglishName: "Georgia",
		DXCC:        []CountryEnum{CountryGeorgia},
	}

	CountryCodeDEU = ISO3166CountryCode{
		Alpha2:      "DE",
		Alpha3:      "DEU",
		Numeric:     "276",
		EnglishName: "Germany",
		DXCC:        []CountryEnum{CountryFederalRepublicOfGermany},
	}

	CountryCodeGHA = ISO3166CountryCode{
		Alpha2:      "GH",
		Alpha3:      "GHA",
		Numeric:     "288",
		EnglishName: "Ghana",
		DXCC:        []CountryEnum{CountryGhana},
	}

	CountryCodeGIB = ISO3166CountryCode{
		Alpha2:      "GI",
		Alpha3:      "GIB",
		Numeric:     "292",
		EnglishName: "Gibraltar",
		DXCC:        []CountryEnum{CountryGibraltar},
	}

	CountryCodeGRC = ISO3166CountryCode{
		Alpha2:      "GR",
		Alpha3:      "GRC",
		Numeric:     "300",
		EnglishName: "Greece",
		DXCC:        []CountryEnum{CountryGreece, CountryCrete, CountryDodecanese, CountryMountAthos},
		Subdivisions: map[string]CountryEnum{
			"M":  CountryCrete,
			"69": CountryMountAthos,
			// Dodecanese forms only a portion of the South Agean administrative region (code L)
		},
	}

	CountryCodeGRL = ISO3166CountryCode{
		Alpha2:      "GL",
		Alpha3:      "GRL",
		Numeric:     "304",
		EnglishName: "Greenland",
		DXCC:        []CountryEnum{CountryGreenland},
	}

	CountryCodeGRD = ISO3166CountryCode{
		Alpha2:      "GD",
		Alpha3:      "GRD",
		Numeric:     "308",
		EnglishName: "Grenada",
		DXCC:        []CountryEnum{CountryGrenada},
	}

	CountryCodeGLP = ISO3166CountryCode{
		Alpha2:      "GP",
		Alpha3:      "GLP",
		Numeric:     "312",
		EnglishName: "Guadeloupe",
		DXCC:        []CountryEnum{CountryGuadeloupe},
	}

	CountryCodeGUM = ISO3166CountryCode{
		Alpha2:      "GU",
		Alpha3:      "GUM",
		Numeric:     "316",
		EnglishName: "Guam",
		DXCC:        []CountryEnum{CountryGuam},
	}

	CountryCodeGTM = ISO3166CountryCode{
		Alpha2:      "GT",
		Alpha3:      "GTM",
		Numeric:     "320",
		EnglishName: "Guatemala",
		DXCC:        []CountryEnum{CountryGuatemala},
	}

	CountryCodeGGY = ISO3166CountryCode{
		Alpha2:      "GG",
		Alpha3:      "GGY",
		Numeric:     "831",
		EnglishName: "Guernsey",
		DXCC:        []CountryEnum{CountryGuernsey},
	}

	CountryCodeGIN = ISO3166CountryCode{
		Alpha2:      "GN",
		Alpha3:      "GIN",
		Numeric:     "324",
		EnglishName: "Guinea",
		DXCC:        []CountryEnum{CountryGuinea},
	}

	CountryCodeGNB = ISO3166CountryCode{
		Alpha2:      "GW",
		Alpha3:      "GNB",
		Numeric:     "624",
		EnglishName: "Guinea-Bissau",
		DXCC:        []CountryEnum{CountryGuineaBissau},
	}

	CountryCodeGUY = ISO3166CountryCode{
		Alpha2:      "GY",
		Alpha3:      "GUY",
		Numeric:     "328",
		EnglishName: "Guyana",
		DXCC:        []CountryEnum{CountryGuyana},
	}

	CountryCodeHTI = ISO3166CountryCode{
		Alpha2:      "HT",
		Alpha3:      "HTI",
		Numeric:     "332",
		EnglishName: "Haiti",
		DXCC:        []CountryEnum{CountryHaiti},
	}

	CountryCodeHMD = ISO3166CountryCode{
		Alpha2:      "HM",
		Alpha3:      "HMD",
		Numeric:     "334",
		EnglishName: "Heard Island and McDonald Islands",
		DXCC:        []CountryEnum{CountryHeardIsland},
	}

	CountryCodeVAT = ISO3166CountryCode{
		Alpha2:      "VA",
		Alpha3:      "VAT",
		Numeric:     "336",
		EnglishName: "Holy See (the)",
		DXCC:        []CountryEnum{CountryVatican},
	}

	CountryCodeHND = ISO3166CountryCode{
		Alpha2:      "HN",
		Alpha3:      "HND",
		Numeric:     "340",
		EnglishName: "Honduras",
		DXCC:        []CountryEnum{CountryHonduras},
	}

	CountryCodeHKG = ISO3166CountryCode{
		Alpha2:      "HK",
		Alpha3:      "HKG",
		Numeric:     "344",
		EnglishName: "Hong Kong",
		DXCC:        []CountryEnum{CountryHongKong},
	}

	CountryCodeHUN = ISO3166CountryCode{
		Alpha2:      "HU",
		Alpha3:      "HUN",
		Numeric:     "348",
		EnglishName: "Hungary",
		DXCC:        []CountryEnum{CountryHungary},
	}

	CountryCodeISL = ISO3166CountryCode{
		Alpha2:      "IS",
		Alpha3:      "ISL",
		Numeric:     "352",
		EnglishName: "Iceland",
		DXCC:        []CountryEnum{CountryIceland},
	}

	CountryCodeIND = ISO3166CountryCode{
		Alpha2:      "IN",
		Alpha3:      "IND",
		Numeric:     "356",
		EnglishName: "India",
		DXCC:        []CountryEnum{CountryIndia, CountryAndamanNicobarIslands, CountryLakshadweepIslands},
		Subdivisions: map[string]CountryEnum{
			"AN": CountryAndamanNicobarIslands,
			"LD": CountryLakshadweepIslands,
		},
	}

	CountryCodeIDN = ISO3166CountryCode{
		Alpha2:      "ID",
		Alpha3:      "IDN",
		Numeric:     "360",
		EnglishName: "Indonesia",
		DXCC:        []CountryEnum{CountryIndonesia},
	}

	CountryCodeIRN = ISO3166CountryCode{
		Alpha2:      "IR",
		Alpha3:      "IRN",
		Numeric:     "364",
		EnglishName: "Iran (Islamic Republic of)",
		DXCC:        []CountryEnum{CountryIran},
	}

	CountryCodeIRQ = ISO3166CountryCode{
		Alpha2:      "IQ",
		Alpha3:      "IRQ",
		Numeric:     "368",
		EnglishName: "Iraq",
		DXCC:        []CountryEnum{CountryIraq},
	}

	CountryCodeIRL = ISO3166CountryCode{
		Alpha2:      "IE",
		Alpha3:      "IRL",
		Numeric:     "372",
		EnglishName: "Ireland",
		DXCC:        []CountryEnum{CountryIreland},
	}

	CountryCodeIMN = ISO3166CountryCode{
		Alpha2:      "IM",
		Alpha3:      "IMN",
		Numeric:     "833",
		EnglishName: "Isle of Man",
		DXCC:        []CountryEnum{CountryIsleOfMan},
	}

	CountryCodeISR = ISO3166CountryCode{
		Alpha2:      "IL",
		Alpha3:      "ISR",
		Numeric:     "376",
		EnglishName: "Israel",
		DXCC:        []CountryEnum{CountryIsrael},
	}

	CountryCodeITA = ISO3166CountryCode{
		Alpha2:      "IT",
		Alpha3:      "ITA",
		Numeric:     "380",
		EnglishName: "Italy",
		DXCC:        []CountryEnum{CountryItaly, CountrySardinia},
		Subdivisions: map[string]CountryEnum{
			"88": CountrySardinia, // Sardinia as an autonomous region
			"CA": CountrySardinia, // Cagliari
			"NU": CountrySardinia, // Nuoro
			"OR": CountrySardinia, // Oristano
			"SS": CountrySardinia, // Sassari
			"SU": CountrySardinia, // Sud Sardegna
			// the following provinces were replaced by Sud Sardegna in 2019
			"CI": CountrySardinia, // Carbonia-Iglesias
			"MD": CountrySardinia, // Medio Campidano non-ISO code
			"OG": CountrySardinia, // Ogliastra
			"OT": CountrySardinia, // Olbia-Tempio
			"VS": CountrySardinia, // ISO code for Medio Campidano (Villacidro & Sanluri)
		},
	}

	CountryCodeJAM = ISO3166CountryCode{
		Alpha2:      "JM",
		Alpha3:      "JAM",
		Numeric:     "388",
		EnglishName: "Jamaica",
		DXCC:        []CountryEnum{CountryJamaica},
	}

	CountryCodeJPN = ISO3166CountryCode{
		Alpha2:      "JP",
		Alpha3:      "JPN",
		Numeric:     "392",
		EnglishName: "Japan",
		DXCC:        []CountryEnum{CountryJapan},
	}

	CountryCodeJEY = ISO3166CountryCode{
		Alpha2:      "JE",
		Alpha3:      "JEY",
		Numeric:     "832",
		EnglishName: "Jersey",
		DXCC:        []CountryEnum{CountryJersey},
	}

	CountryCodeJOR = ISO3166CountryCode{
		Alpha2:      "JO",
		Alpha3:      "JOR",
		Numeric:     "400",
		EnglishName: "Jordan",
		DXCC:        []CountryEnum{CountryJordan},
	}

	CountryCodeKAZ = ISO3166CountryCode{
		Alpha2:      "KZ",
		Alpha3:      "KAZ",
		Numeric:     "398",
		EnglishName: "Kazakhstan",
		DXCC:        []CountryEnum{CountryKazakhstan},
	}

	CountryCodeKEN = ISO3166CountryCode{
		Alpha2:      "KE",
		Alpha3:      "KEN",
		Numeric:     "404",
		EnglishName: "Kenya",
		DXCC:        []CountryEnum{CountryKenya},
	}

	CountryCodeKIR = ISO3166CountryCode{
		Alpha2:      "KI",
		Alpha3:      "KIR",
		Numeric:     "296",
		EnglishName: "Kiribati",
		DXCC:        []CountryEnum{CountryWKiribatiGilbertIslands, CountryEKiribatiLineIslands, CountryCKiribatiBritishPhoenixIslands, CountryBanabaIslandOceanIsland},
		Subdivisions: map[string]CountryEnum{
			"G": CountryWKiribatiGilbertIslands,
			"L": CountryEKiribatiLineIslands,
			"P": CountryCKiribatiBritishPhoenixIslands,
			// Banaba Island doesn't have its own ISO code
		},
	}

	CountryCodePRK = ISO3166CountryCode{
		Alpha2:      "KP",
		Alpha3:      "PRK",
		Numeric:     "408",
		EnglishName: "Korea (the Democratic People's Republic of)",
		DXCC:        []CountryEnum{CountryDemocraticPeoplesRepOfKorea},
	}

	CountryCodeKOR = ISO3166CountryCode{
		Alpha2:      "KR",
		Alpha3:      "KOR",
		Numeric:     "410",
		EnglishName: "Korea (the Republic of)",
		DXCC:        []CountryEnum{CountryRepublicOfKorea},
	}

	CountryCodeKWT = ISO3166CountryCode{
		Alpha2:      "KW",
		Alpha3:      "KWT",
		Numeric:     "414",
		EnglishName: "Kuwait",
		DXCC:        []CountryEnum{CountryKuwait},
	}

	CountryCodeKGZ = ISO3166CountryCode{
		Alpha2:      "KG",
		Alpha3:      "KGZ",
		Numeric:     "417",
		EnglishName: "Kyrgyzstan",
		DXCC:        []CountryEnum{CountryKyrgyzstan},
	}

	CountryCodeLAO = ISO3166CountryCode{
		Alpha2:      "LA",
		Alpha3:      "LAO",
		Numeric:     "418",
		EnglishName: "Lao People's Democratic Republic (the)",
		DXCC:        []CountryEnum{CountryLaos},
	}

	CountryCodeLVA = ISO3166CountryCode{
		Alpha2:      "LV",
		Alpha3:      "LVA",
		Numeric:     "428",
		EnglishName: "Latvia",
		DXCC:        []CountryEnum{CountryLatvia},
	}

	CountryCodeLBN = ISO3166CountryCode{
		Alpha2:      "LB",
		Alpha3:      "LBN",
		Numeric:     "422",
		EnglishName: "Lebanon",
		DXCC:        []CountryEnum{CountryLebanon},
	}

	CountryCodeLSO = ISO3166CountryCode{
		Alpha2:      "LS",
		Alpha3:      "LSO",
		Numeric:     "426",
		EnglishName: "Lesotho",
		DXCC:        []CountryEnum{CountryLesotho},
	}

	CountryCodeLBR = ISO3166CountryCode{
		Alpha2:      "LR",
		Alpha3:      "LBR",
		Numeric:     "430",
		EnglishName: "Liberia",
		DXCC:        []CountryEnum{CountryLiberia},
	}

	CountryCodeLBY = ISO3166CountryCode{
		Alpha2:      "LY",
		Alpha3:      "LBY",
		Numeric:     "434",
		EnglishName: "Libya",
		DXCC:        []CountryEnum{CountryLibya},
	}

	CountryCodeLIE = ISO3166CountryCode{
		Alpha2:      "LI",
		Alpha3:      "LIE",
		Numeric:     "438",
		EnglishName: "Liechtenstein",
		DXCC:        []CountryEnum{CountryLiechtenstein},
	}

	CountryCodeLTU = ISO3166CountryCode{
		Alpha2:      "LT",
		Alpha3:      "LTU",
		Numeric:     "440",
		EnglishName: "Lithuania",
		DXCC:        []CountryEnum{CountryLithuania},
	}

	CountryCodeLUX = ISO3166CountryCode{
		Alpha2:      "LU",
		Alpha3:      "LUX",
		Numeric:     "442",
		EnglishName: "Luxembourg",
		DXCC:        []CountryEnum{CountryLuxembourg},
	}

	CountryCodeMAC = ISO3166CountryCode{
		Alpha2:      "MO",
		Alpha3:      "MAC",
		Numeric:     "446",
		EnglishName: "Macao",
		DXCC:        []CountryEnum{CountryMacao},
	}

	CountryCodeMDG = ISO3166CountryCode{
		Alpha2:      "MG",
		Alpha3:      "MDG",
		Numeric:     "450",
		EnglishName: "Madagascar",
		DXCC:        []CountryEnum{CountryMadagascar},
	}

	CountryCodeMWI = ISO3166CountryCode{
		Alpha2:      "MW",
		Alpha3:      "MWI",
		Numeric:     "454",
		EnglishName: "Malawi",
		DXCC:        []CountryEnum{CountryMalawi},
	}

	CountryCodeMYS = ISO3166CountryCode{
		Alpha2:      "MY",
		Alpha3:      "MYS",
		Numeric:     "458",
		EnglishName: "Malaysia",
		DXCC:        []CountryEnum{CountryWestMalaysia, CountryEastMalaysia},
		Subdivisions: map[string]CountryEnum{
			"12": CountryEastMalaysia, // Sabah
			"13": CountryEastMalaysia, // Sarawak
			"15": CountryEastMalaysia, // Labuan
			"01": CountryWestMalaysia, // Johor
			"02": CountryWestMalaysia, // Kedah
			"03": CountryWestMalaysia, // Kelantan
			"04": CountryWestMalaysia, // Melaka
			"05": CountryWestMalaysia, // Negeri Sembilan
			"06": CountryWestMalaysia, // Pahang
			"07": CountryWestMalaysia, // Pulau Pinang
			"08": CountryWestMalaysia, // Perak
			"09": CountryWestMalaysia, // Perlis
			"10": CountryWestMalaysia, // Selangor
			"11": CountryWestMalaysia, // Terengganu
			"14": CountryWestMalaysia, // Kuala Lumpur
			"16": CountryWestMalaysia, // Putrajaya
		},
	}

	CountryCodeMDV = ISO3166CountryCode{
		Alpha2:      "MV",
		Alpha3:      "MDV",
		Numeric:     "462",
		EnglishName: "Maldives",
		DXCC:        []CountryEnum{CountryMaldives},
	}

	CountryCodeMLI = ISO3166CountryCode{
		Alpha2:      "ML",
		Alpha3:      "MLI",
		Numeric:     "466",
		EnglishName: "Mali",
		DXCC:        []CountryEnum{CountryMali},
	}

	CountryCodeMLT = ISO3166CountryCode{
		Alpha2:      "MT",
		Alpha3:      "MLT",
		Numeric:     "470",
		EnglishName: "Malta",
		DXCC:        []CountryEnum{CountryMalta},
	}

	CountryCodeMHL = ISO3166CountryCode{
		Alpha2:      "MH",
		Alpha3:      "MHL",
		Numeric:     "584",
		EnglishName: "Marshall Islands (the)",
		DXCC:        []CountryEnum{CountryMarshallIslands},
	}

	CountryCodeMTQ = ISO3166CountryCode{
		Alpha2:      "MQ",
		Alpha3:      "MTQ",
		Numeric:     "474",
		EnglishName: "Martinique",
		DXCC:        []CountryEnum{CountryMartinique},
	}

	CountryCodeMRT = ISO3166CountryCode{
		Alpha2:      "MR",
		Alpha3:      "MRT",
		Numeric:     "478",
		EnglishName: "Mauritania",
		DXCC:        []CountryEnum{CountryMauritania},
	}

	CountryCodeMUS = ISO3166CountryCode{
		Alpha2:      "MU",
		Alpha3:      "MUS",
		Numeric:     "480",
		EnglishName: "Mauritius",
		DXCC:        []CountryEnum{CountryMauritius, CountryAgalegaStBrandonIslands, CountryRodriguesIsland},
		Subdivisions: map[string]CountryEnum{
			"AG": CountryAgalegaStBrandonIslands, // Agalega Islands
			"CC": CountryAgalegaStBrandonIslands, // Cargados Carajos Shoals
			"RO": CountryRodriguesIsland,
		},
	}

	CountryCodeMYT = ISO3166CountryCode{
		Alpha2:      "YT",
		Alpha3:      "MYT",
		Numeric:     "175",
		EnglishName: "Mayotte",
		DXCC:        []CountryEnum{CountryMayotte},
	}

	CountryCodeMEX = ISO3166CountryCode{
		Alpha2:      "MX",
		Alpha3:      "MEX",
		Numeric:     "484",
		EnglishName: "Mexico",
		DXCC:        []CountryEnum{CountryMexico},
	}

	CountryCodeFSM = ISO3166CountryCode{
		Alpha2:      "FM",
		Alpha3:      "FSM",
		Numeric:     "583",
		EnglishName: "Micronesia (Federated States of)",
		DXCC:        []CountryEnum{CountryMicronesia},
	}

	CountryCodeMDA = ISO3166CountryCode{
		Alpha2:      "MD",
		Alpha3:      "MDA",
		Numeric:     "498",
		EnglishName: "Moldova (the Republic of)",
		DXCC:        []CountryEnum{CountryMoldova},
	}

	CountryCodeMCO = ISO3166CountryCode{
		Alpha2:      "MC",
		Alpha3:      "MCO",
		Numeric:     "492",
		EnglishName: "Monaco",
		DXCC:        []CountryEnum{CountryMonaco},
	}

	CountryCodeMNG = ISO3166CountryCode{
		Alpha2:      "MN",
		Alpha3:      "MNG",
		Numeric:     "496",
		EnglishName: "Mongolia",
		DXCC:        []CountryEnum{CountryMongolia},
	}

	CountryCodeMNE = ISO3166CountryCode{
		Alpha2:      "ME",
		Alpha3:      "MNE",
		Numeric:     "499",
		EnglishName: "Montenegro",
		DXCC:        []CountryEnum{CountryMontenegro},
	}

	CountryCodeMSR = ISO3166CountryCode{
		Alpha2:      "MS",
		Alpha3:      "MSR",
		Numeric:     "500",
		EnglishName: "Montserrat",
		DXCC:        []CountryEnum{CountryMontserrat},
	}

	CountryCodeMAR = ISO3166CountryCode{
		Alpha2:      "MA",
		Alpha3:      "MAR",
		Numeric:     "504",
		EnglishName: "Morocco",
		DXCC:        []CountryEnum{CountryMorocco},
	}

	CountryCodeMOZ = ISO3166CountryCode{
		Alpha2:      "MZ",
		Alpha3:      "MOZ",
		Numeric:     "508",
		EnglishName: "Mozambique",
		DXCC:        []CountryEnum{CountryMozambique},
	}

	CountryCodeMMR = ISO3166CountryCode{
		Alpha2:      "MM",
		Alpha3:      "MMR",
		Numeric:     "104",
		EnglishName: "Myanmar",
		DXCC:        []CountryEnum{CountryMyanmar},
	}

	CountryCodeNAM = ISO3166CountryCode{
		Alpha2:      "NA",
		Alpha3:      "NAM",
		Numeric:     "516",
		EnglishName: "Namibia",
		DXCC:        []CountryEnum{CountryNamibia},
	}

	CountryCodeNRU = ISO3166CountryCode{
		Alpha2:      "NR",
		Alpha3:      "NRU",
		Numeric:     "520",
		EnglishName: "Nauru",
		DXCC:        []CountryEnum{CountryNauru},
	}

	CountryCodeNPL = ISO3166CountryCode{
		Alpha2:      "NP",
		Alpha3:      "NPL",
		Numeric:     "524",
		EnglishName: "Nepal",
		DXCC:        []CountryEnum{CountryNepal},
	}

	CountryCodeNLD = ISO3166CountryCode{
		Alpha2:      "NL",
		Alpha3:      "NLD",
		Numeric:     "528",
		EnglishName: "Netherlands (the)",
		DXCC:        []CountryEnum{CountryNetherlands},
	}

	CountryCodeNCL = ISO3166CountryCode{
		Alpha2:      "NC",
		Alpha3:      "NCL",
		Numeric:     "540",
		EnglishName: "New Caledonia",
		DXCC:        []CountryEnum{CountryNewCaledonia},
	}

	CountryCodeNZL = ISO3166CountryCode{
		Alpha2:      "NZ",
		Alpha3:      "NZL",
		Numeric:     "554",
		EnglishName: "New Zealand",
		DXCC:        []CountryEnum{CountryNewZealand, CountryChathamIslands},
		Subdivisions: map[string]CountryEnum{
			// Chatham and other New Zealand territories don't have subdivision codes, so regions
			// from North Island and South Island are listed here to ensure contacts logged with "NZ"
			// are interpreted as being from New Zealand
			"AUK": CountryNewZealand, // Auckland
			"BOP": CountryNewZealand, // Bay of Plenty
			"CAN": CountryNewZealand, // Canterbury
			"GIS": CountryNewZealand, // Gisborne
			"HKB": CountryNewZealand, // Hawke's Bay
			"MBH": CountryNewZealand, // Marlborough
			"MWT": CountryNewZealand, // Manawatu-Whanganui
			"NSN": CountryNewZealand, // Nelson
			"NTL": CountryNewZealand, // Northland
			"OTA": CountryNewZealand, // Otago
			"STL": CountryNewZealand, // Southland
			"TAS": CountryNewZealand, // Tasman
			"TKI": CountryNewZealand, // Taranaki
			"WGN": CountryNewZealand, // Greater Wellington
			"WKO": CountryNewZealand, // Waikato
			"WTC": CountryNewZealand, // West Coast
		},
	}

	CountryCodeNIC = ISO3166CountryCode{
		Alpha2:      "NI",
		Alpha3:      "NIC",
		Numeric:     "558",
		EnglishName: "Nicaragua",
		DXCC:        []CountryEnum{CountryNicaragua},
	}

	CountryCodeNER = ISO3166CountryCode{
		Alpha2:      "NE",
		Alpha3:      "NER",
		Numeric:     "562",
		EnglishName: "Niger (the)",
		DXCC:        []CountryEnum{CountryNiger},
	}

	CountryCodeNGA = ISO3166CountryCode{
		Alpha2:      "NG",
		Alpha3:      "NGA",
		Numeric:     "566",
		EnglishName: "Nigeria",
		DXCC:        []CountryEnum{CountryNigeria},
	}

	CountryCodeNIU = ISO3166CountryCode{
		Alpha2:      "NU",
		Alpha3:      "NIU",
		Numeric:     "570",
		EnglishName: "Niue",
		DXCC:        []CountryEnum{CountryNiue},
	}

	CountryCodeNFK = ISO3166CountryCode{
		Alpha2:      "NF",
		Alpha3:      "NFK",
		Numeric:     "574",
		EnglishName: "Norfolk Island",
		DXCC:        []CountryEnum{CountryNorfolkIsland},
	}

	CountryCodeMNP = ISO3166CountryCode{
		Alpha2:      "MP",
		Alpha3:      "MNP",
		Numeric:     "580",
		EnglishName: "Northern Mariana Islands (the)",
		DXCC:        []CountryEnum{CountryMarianaIslands},
	}

	CountryCodeNOR = ISO3166CountryCode{
		Alpha2:      "NO",
		Alpha3:      "NOR",
		Numeric:     "578",
		EnglishName: "Norway",
		DXCC:        []CountryEnum{CountryNorway},
	}

	CountryCodeOMN = ISO3166CountryCode{
		Alpha2:      "OM",
		Alpha3:      "OMN",
		Numeric:     "512",
		EnglishName: "Oman",
		DXCC:        []CountryEnum{CountryOman},
	}

	CountryCodePAK = ISO3166CountryCode{
		Alpha2:      "PK",
		Alpha3:      "PAK",
		Numeric:     "586",
		EnglishName: "Pakistan",
		DXCC:        []CountryEnum{CountryPakistan},
	}

	CountryCodePLW = ISO3166CountryCode{
		Alpha2:      "PW",
		Alpha3:      "PLW",
		Numeric:     "585",
		EnglishName: "Palau",
		DXCC:        []CountryEnum{CountryPalau},
	}

	CountryCodePSE = ISO3166CountryCode{
		Alpha2:      "PS",
		Alpha3:      "PSE",
		Numeric:     "275",
		EnglishName: "Palestine, State of",
		DXCC:        []CountryEnum{CountryPalestine},
	}

	CountryCodePAN = ISO3166CountryCode{
		Alpha2:      "PA",
		Alpha3:      "PAN",
		Numeric:     "591",
		EnglishName: "Panama",
		DXCC:        []CountryEnum{CountryPanama},
	}

	CountryCodePNG = ISO3166CountryCode{
		Alpha2:      "PG",
		Alpha3:      "PNG",
		Numeric:     "598",
		EnglishName: "Papua New Guinea",
		DXCC:        []CountryEnum{CountryPapuaNewGuinea},
	}

	CountryCodePRY = ISO3166CountryCode{
		Alpha2:      "PY",
		Alpha3:      "PRY",
		Numeric:     "600",
		EnglishName: "Paraguay",
		DXCC:        []CountryEnum{CountryParaguay},
	}

	CountryCodePER = ISO3166CountryCode{
		Alpha2:      "PE",
		Alpha3:      "PER",
		Numeric:     "604",
		EnglishName: "Peru",
		DXCC:        []CountryEnum{CountryPeru},
	}

	CountryCodePHL = ISO3166CountryCode{
		Alpha2:      "PH",
		Alpha3:      "PHL",
		Numeric:     "608",
		EnglishName: "Philippines (the)",
		DXCC:        []CountryEnum{CountryPhilippines},
	}

	CountryCodePCN = ISO3166CountryCode{
		Alpha2:      "PN",
		Alpha3:      "PCN",
		Numeric:     "612",
		EnglishName: "Pitcairn",
		DXCC:        []CountryEnum{CountryPitcairnIsland},
	}

	CountryCodePOL = ISO3166CountryCode{
		Alpha2:      "PL",
		Alpha3:      "POL",
		Numeric:     "616",
		EnglishName: "Poland",
		DXCC:        []CountryEnum{CountryPoland},
	}

	CountryCodePRT = ISO3166CountryCode{
		Alpha2:      "PT",
		Alpha3:      "PRT",
		Numeric:     "620",
		EnglishName: "Portugal",
		DXCC:        []CountryEnum{CountryPortugal, CountryAzores, CountryMadeiraIslands},
		Subdivisions: map[string]CountryEnum{
			"20": CountryAzores,         // ISO code
			"AC": CountryAzores,         // ADIF code
			"30": CountryMadeiraIslands, // ISO code
			"MD": CountryMadeiraIslands, // ADIF code
		},
	}

	CountryCodePRI = ISO3166CountryCode{
		Alpha2:      "PR",
		Alpha3:      "PRI",
		Numeric:     "630",
		EnglishName: "Puerto Rico",
		DXCC:        []CountryEnum{CountryPuertoRico},
	}

	CountryCodeQAT = ISO3166CountryCode{
		Alpha2:      "QA",
		Alpha3:      "QAT",
		Numeric:     "634",
		EnglishName: "Qatar",
		DXCC:        []CountryEnum{CountryQatar},
	}

	CountryCodeMKD = ISO3166CountryCode{
		Alpha2:      "MK",
		Alpha3:      "MKD",
		Numeric:     "807",
		EnglishName: "Republic of North Macedonia",
		DXCC:        []CountryEnum{CountryNorthMacedoniaRepublicOf},
	}

	CountryCodeROU = ISO3166CountryCode{
		Alpha2:      "RO",
		Alpha3:      "ROU",
		Numeric:     "642",
		EnglishName: "Romania",
		DXCC:        []CountryEnum{CountryRomania},
	}

	CountryCodeRUS = ISO3166CountryCode{
		Alpha2:      "RU",
		Alpha3:      "RUS",
		Numeric:     "643",
		EnglishName: "Russian Federation (the)",
		DXCC:        []CountryEnum{CountryEuropeanRussia, CountryAsiaticRussia, CountryKaliningrad},
		Subdivisions: map[string]CountryEnum{
			// Asia side
			"AL":  CountryAsiaticRussia, // Altay Republic (ISO), Altaysky Kraj (ADIF)
			"ALT": CountryAsiaticRussia, // Altayskiy Kraj (ISO)
			"AM":  CountryAsiaticRussia, // Amurskaya Oblast (ADIF)
			"AMU": CountryAsiaticRussia, // Amurskaya Oblast (ISO)
			"BA":  CountryAsiaticRussia, // Bashkortostan Republic (ISO, ADIF)
			"BU":  CountryAsiaticRussia, // Buryatiya Republic (ISO, ADIF)
			"CB":  CountryAsiaticRussia, // Chelyabinskaya Oblast (ADIF)
			"CHE": CountryAsiaticRussia, // Chelyabinskaya Oblast (ISO)
			"CHU": CountryAsiaticRussia, // Chukotskiy Autonomous Okrug (ISO)
			"CK":  CountryAsiaticRussia, // Chukotskiy Autonomous Okrug (ADIF)
			"CT":  CountryAsiaticRussia, // Zabaykalsky Kraj (ADIF)
			"EA":  CountryAsiaticRussia, // Yevreyskaya Autonomous Oblast (ADIF)
			"GA":  CountryAsiaticRussia, // Altaj Respublika (ADIF)
			"HA":  CountryAsiaticRussia, // Hakasija Respublika (ADIF)
			"HK":  CountryAsiaticRussia, // Khabarovskiy Kraj (ADIF)
			"HM":  CountryAsiaticRussia, // Khanty-Mansiyskiy Autonomous Okrug (ADIF)
			"IR":  CountryAsiaticRussia, // Irkutskaya Oblast (ADIF)
			"IRK": CountryAsiaticRussia, // Irkutskaya Oblast (ISO)
			"KAM": CountryAsiaticRussia, // Kamchatskiy Kray (ISO)
			"KE":  CountryAsiaticRussia, // Kemerovskaya Oblast (ADIF)
			"KEM": CountryAsiaticRussia, // Kemerovskaya Oblast (ISO)
			"KGN": CountryAsiaticRussia, // Kurganskaya Oblast (ISO)
			"KHA": CountryAsiaticRussia, // Khabarovskiy Kraj (ISO)
			"KHM": CountryAsiaticRussia, // Khanty-Mansiyskiy Autonomous Okrug (ISO)
			"KK":  CountryAsiaticRussia, // Khaskasiya Republic (ISO), Krasnoyarsk Kraj (ADIF)
			"KN":  CountryAsiaticRussia, // Kurganskaya Oblast (ADIF)
			"KO":  CountryAsiaticRussia, // Komi Republic (ISO, ADIF)
			"KT":  CountryAsiaticRussia, // Kamchatskaya Oblast (ADIF)
			"KYA": CountryAsiaticRussia, // Krasnojarskij Kraj (ADIF)
			"MAG": CountryAsiaticRussia, // Magadanskaya Oblast (ISO)
			"MG":  CountryAsiaticRussia, // Magadanskaya Oblast (ADIF)
			"NS":  CountryAsiaticRussia, // Novosibriskaya Oblast (ADIF)
			"NVS": CountryAsiaticRussia, // Novosibriskaya Oblast (ISO)
			"OB":  CountryAsiaticRussia, // Orenburgskaya Oblast (ADIF)
			"OM":  CountryAsiaticRussia, // Omskaya Oblast (ADIF)
			"OMS": CountryAsiaticRussia, // Omskaya Oblast (ISO)
			"ORE": CountryAsiaticRussia, // Orenburgskaya Oblast (ISO)
			"PER": CountryAsiaticRussia, // Permskiy Kraj (ISO)
			"PK":  CountryAsiaticRussia, // Primorsky Kraj (ADIF)
			"PM":  CountryAsiaticRussia, // Permskiy Kraj (ADIF)
			"PRI": CountryAsiaticRussia, // Primorsky Kraj (ISO)
			// "SA": CountryAsiaticRussia, // Sakha Republic (ISO), conflicts with Europe Saratovskaya (ADIF)
			"SAK": CountryAsiaticRussia, // Sakhalinskaya Oblast (ISO)
			"SL":  CountryAsiaticRussia, // Sakhalinskaya Oblast (ADIF)
			"SV":  CountryAsiaticRussia, // Sverdlovskaya Oblast (ADIF)
			"SVE": CountryAsiaticRussia, // Sverdlovskaya Oblast (ISO)
			"TN":  CountryAsiaticRussia, // Tyumenskaya Oblast (ADIF)
			"TO":  CountryAsiaticRussia, // Tomskaya Oblast (ADIF)
			"TOM": CountryAsiaticRussia, // Tomskaya Oblast (ISO)
			"TU":  CountryAsiaticRussia, // Tuva Republic (ADIF)
			"TY":  CountryAsiaticRussia, // Tuva Republic (ISO)
			"TYU": CountryAsiaticRussia, // Tyumenskaya Republic (ISO)
			"YA":  CountryAsiaticRussia, // Sakha Republic (ISO)
			"YAN": CountryAsiaticRussia, // Yamalo-Nenetsky Autonomous Oblast (ISO)
			"YEV": CountryAsiaticRussia, // Yevreyskaya Autonomous Oblast (ISO)
			"YN":  CountryAsiaticRussia, // Yamalo-Nenetsky Autonomous Oblast (ADIF)
			"ZAB": CountryAsiaticRussia, // Zabaykalsky Kraj (ISO)
			// Europe side
			"AD":  CountryEuropeanRussia, // Adygeya Republic (ISO, ADIF)
			"AO":  CountryEuropeanRussia, // Astrakhanskaya Oblast (ADIF)
			"AR":  CountryEuropeanRussia, // Arkhangelskaya Oblast (ADIF)
			"ARK": CountryEuropeanRussia, // Arkhangelskaya Oblast (ISO)
			"AST": CountryEuropeanRussia, // Astrakhanskaya Oblast (ISO)
			"BEL": CountryEuropeanRussia, // Belgorodskaya Oblast (ISO)
			"BO":  CountryEuropeanRussia, // Belgorodskaya Oblast (ADIF)
			"BR":  CountryEuropeanRussia, // Bryanskaya Oblast (ADIF)
			"BRY": CountryEuropeanRussia, // Bryanskaya Oblast (ISO)
			"CE":  CountryEuropeanRussia, // Chechnya Republic (ISO)
			"CN":  CountryEuropeanRussia, // Chechnya Republic (ADIF)
			"CU":  CountryEuropeanRussia, // Chuvashia Republic (ISO, ADIF)
			"DA":  CountryEuropeanRussia, // Daghestan Republic (ISO, ADIF)
			"IN":  CountryEuropeanRussia, // Ingushetia Republic (ISO, ADIF)
			"IV":  CountryEuropeanRussia, // Ivanovskaya Oblast (ADIF)
			"IVA": CountryEuropeanRussia, // Ivanovskaya Oblast (ADIF)
			"KB":  CountryEuropeanRussia, // Kabardino-Balkaria Republic (ISO, ADIF)
			"KDA": CountryEuropeanRussia, // Krasnodarskiy Kraj (ISO)
			"KG":  CountryEuropeanRussia, // Kaluzhskaya Oblast (ADIF)
			"KI":  CountryEuropeanRussia, // Kirovskaya Oblast (ADIF)
			"KIR": CountryEuropeanRussia, // Kirovskaya Oblast (ADIF)
			"KL":  CountryEuropeanRussia, // Kalmykia Republic (ISO), Karelia Republic (ADIF)
			"KLU": CountryEuropeanRussia, // Kaluzhskaya Oblast (ISO)
			"KM":  CountryEuropeanRussia, // Kalmykia Republic (ADIF)
			"KOS": CountryEuropeanRussia, // Kostromskaya Oblast (ADIF)
			"KR":  CountryEuropeanRussia, // Karelia Republic (ISO)
			"KRS": CountryEuropeanRussia, // Kurskaya Oblast (ISO)
			"KS":  CountryEuropeanRussia, // Kostromskaya Oblast (ADIF)
			"KU":  CountryEuropeanRussia, // Kurskaya Oblast (ADIF)
			"LEN": CountryEuropeanRussia, // Leningraskaya Oblast (ISO)
			"LIP": CountryEuropeanRussia, // Lipetskaya Oblast (ISO)
			"LO":  CountryEuropeanRussia, // Leningraskaya Oblast (ADIF)
			"LP":  CountryEuropeanRussia, // Lipetskaya Oblast (ADIF)
			"MA":  CountryEuropeanRussia, // Moskow City (ADIF)
			"MD":  CountryEuropeanRussia, // Mordovia Republic (ADIF)
			"ME":  CountryEuropeanRussia, // Marij-El Republic (ISO)
			"MO":  CountryEuropeanRussia, // Mordovia Republic (ISO), Moscowskaya Oblast (ADIF)
			"MOS": CountryEuropeanRussia, // Moscowskaya Oblast (ISO)
			"MOW": CountryEuropeanRussia, // Moscow City (ISO)
			"MR":  CountryEuropeanRussia, // Marij-El Republic (ADIF)
			"MU":  CountryEuropeanRussia, // Murmanskaya Oblast (ADIF)
			"MUR": CountryEuropeanRussia, // Murmanskaya Oblast (ISO)
			"NEN": CountryEuropeanRussia, // Nenetsky Autonomous Okrug (ISO)
			"NGR": CountryEuropeanRussia, // Novgoroskaya Oblast (ISO)
			"NN":  CountryEuropeanRussia, // Nizhegorodskaya Oblast (ADIF)
			"NO":  CountryEuropeanRussia, // Nenetsky Autonomous Okrug (ADIF)
			"NV":  CountryEuropeanRussia, // Novgoroskaya Oblast (ADIF)
			"OR":  CountryEuropeanRussia, // Orlovskaya Oblast (ADIF)
			"ORL": CountryEuropeanRussia, // Orlovskaya Oblast (ISO)
			"PE":  CountryEuropeanRussia, // Penzenskaya Oblast (ADIF)
			"PNZ": CountryEuropeanRussia, // Penzenskaya Oblast (ISO)
			"PS":  CountryEuropeanRussia, // Pskovskaya Oblast (ADIF)
			"PSK": CountryEuropeanRussia, // Pskovskaya Oblast (ISO)
			"RA":  CountryEuropeanRussia, // Ryazanskaya Oblast (ADIF)
			"RO":  CountryEuropeanRussia, // Rostovskaya Oblast (ADIF)
			"ROS": CountryEuropeanRussia, // Rostovskaya Oblast (ISO)
			"RYA": CountryEuropeanRussia, // Ryazanskaya Oblast (ISO)
			// "SA": CountryEuropeanRussia, // Saratovskaya Oblast (ADIF), conflicts with Asia Sakha (ISO)
			"SAM": CountryEuropeanRussia, // Samarskaya Oblast (ADIF)
			"SAR": CountryEuropeanRussia, // Saratovskaya Oblast (ISO)
			"SE":  CountryEuropeanRussia, // Northern Ossetia Republic (ISO)
			"SM":  CountryEuropeanRussia, // Smolenskaya Oblast (ADIF)
			"SMO": CountryEuropeanRussia, // Smolenskaya Oblast (ISO)
			"SO":  CountryEuropeanRussia, // Northern Ossetia Republic (ADIF)
			"SP":  CountryEuropeanRussia, // St. Petersburg City (ADIF)
			"SPE": CountryEuropeanRussia, // St. Petersburg City (ADIF)
			"SR":  CountryEuropeanRussia, // Samaraskaya Oblast (ISO)
			"ST":  CountryEuropeanRussia, // Stavropolsky Kraj (ADIF)
			"STA": CountryEuropeanRussia, // Stavropolsky Kraj (ISO)
			"TA":  CountryEuropeanRussia, // Tataria Republic (ISO, ADIF)
			"TAM": CountryEuropeanRussia, // Tambovskaya Oblast (ISO)
			"TB":  CountryEuropeanRussia, // Tambovskaya Oblast (ADIF)
			"TL":  CountryEuropeanRussia, // Tulskaya Oblast (ADIF)
			"TUL": CountryEuropeanRussia, // Tulskaya Oblast (ISO)
			"TV":  CountryEuropeanRussia, // Tverskaya Oblast (ISO)
			"TVE": CountryEuropeanRussia, // Tverskaya Oblast (ISO)
			"UD":  CountryEuropeanRussia, // Udmurtia Republic (ISO, ADIF)
			"UL":  CountryEuropeanRussia, // Ulyanovskaya Oblast (ADIF)
			"ULY": CountryEuropeanRussia, // Ulyanovskaya Oblast (ISO)
			"VG":  CountryEuropeanRussia, // Volgogradskaya Oblast (ADIF)
			"VGG": CountryEuropeanRussia, // Volgogradskaya Oblast (ISO)
			"VL":  CountryEuropeanRussia, // Vladimirskaya Oblast (ADIF)
			"VLA": CountryEuropeanRussia, // Vladimirskaya Oblast (ISO)
			"VLG": CountryEuropeanRussia, // Vologodskaya Oblast (ISO)
			"VO":  CountryEuropeanRussia, // Vologodskaya Oblast (ADIF)
			"VOR": CountryEuropeanRussia, // Voronezhskaya Oblast (ISO)
			"VR":  CountryEuropeanRussia, // Voronezhskaya Oblast (ADIF)
			"YAR": CountryEuropeanRussia, // Yaroslavskaya Oblast (ISO)
			"YR":  CountryEuropeanRussia, // Yaroslavskaya Oblast (ADIF)
			// Kaliningrad
			"KA":  CountryKaliningrad, // Kaliningraskaya Oblast (ADIF)
			"KGD": CountryKaliningrad, // Kaliningraskaya Oblast (ISO)
		},
	}

	CountryCodeRWA = ISO3166CountryCode{
		Alpha2:      "RW",
		Alpha3:      "RWA",
		Numeric:     "646",
		EnglishName: "Rwanda",
		DXCC:        []CountryEnum{CountryRwanda},
	}

	CountryCodeREU = ISO3166CountryCode{
		Alpha2:      "RE",
		Alpha3:      "REU",
		Numeric:     "638",
		EnglishName: "Réunion",
		DXCC:        []CountryEnum{CountryReunionIsland},
	}

	CountryCodeBLM = ISO3166CountryCode{
		Alpha2:      "BL",
		Alpha3:      "BLM",
		Numeric:     "652",
		EnglishName: "Saint Barthélemy",
		DXCC:        []CountryEnum{CountrySaintBarthelemy},
	}

	CountryCodeSHN = ISO3166CountryCode{
		Alpha2:      "SH",
		Alpha3:      "SHN",
		Numeric:     "654",
		EnglishName: "Saint Helena, Ascension and Tristan da Cunha",
		DXCC:        []CountryEnum{CountryStHelena, CountryAscensionIsland, CountryTristanDaCunhaGoughIsland},
		Subdivisions: map[string]CountryEnum{
			"AC": CountryAscensionIsland,
			"HL": CountryStHelena,
			"TA": CountryTristanDaCunhaGoughIsland,
		},
	}

	CountryCodeKNA = ISO3166CountryCode{
		Alpha2:      "KN",
		Alpha3:      "KNA",
		Numeric:     "659",
		EnglishName: "Saint Kitts and Nevis",
		DXCC:        []CountryEnum{CountryStKittsNevis},
	}

	CountryCodeLCA = ISO3166CountryCode{
		Alpha2:      "LC",
		Alpha3:      "LCA",
		Numeric:     "662",
		EnglishName: "Saint Lucia",
		DXCC:        []CountryEnum{CountryStLucia},
	}

	CountryCodeMAF = ISO3166CountryCode{
		Alpha2:      "MF",
		Alpha3:      "MAF",
		Numeric:     "663",
		EnglishName: "Saint Martin (French part)",
		DXCC:        []CountryEnum{CountrySaintMartin},
	}

	CountryCodeSPM = ISO3166CountryCode{
		Alpha2:      "PM",
		Alpha3:      "SPM",
		Numeric:     "666",
		EnglishName: "Saint Pierre and Miquelon",
		DXCC:        []CountryEnum{CountryStPierreMiquelon},
	}

	CountryCodeVCT = ISO3166CountryCode{
		Alpha2:      "VC",
		Alpha3:      "VCT",
		Numeric:     "670",
		EnglishName: "Saint Vincent and the Grenadines",
		DXCC:        []CountryEnum{CountryStVincent},
	}

	CountryCodeWSM = ISO3166CountryCode{
		Alpha2:      "WS",
		Alpha3:      "WSM",
		Numeric:     "882",
		EnglishName: "Samoa",
		DXCC:        []CountryEnum{CountrySamoa},
	}

	CountryCodeSMR = ISO3166CountryCode{
		Alpha2:      "SM",
		Alpha3:      "SMR",
		Numeric:     "674",
		EnglishName: "San Marino",
		DXCC:        []CountryEnum{CountrySanMarino},
	}

	CountryCodeSTP = ISO3166CountryCode{
		Alpha2:      "ST",
		Alpha3:      "STP",
		Numeric:     "678",
		EnglishName: "Sao Tome and Principe",
		DXCC:        []CountryEnum{CountrySaoTomePrincipe},
	}

	CountryCodeSAU = ISO3166CountryCode{
		Alpha2:      "SA",
		Alpha3:      "SAU",
		Numeric:     "682",
		EnglishName: "Saudi Arabia",
		DXCC:        []CountryEnum{CountrySaudiArabia},
	}

	CountryCodeSEN = ISO3166CountryCode{
		Alpha2:      "SN",
		Alpha3:      "SEN",
		Numeric:     "686",
		EnglishName: "Senegal",
		DXCC:        []CountryEnum{CountrySenegal},
	}

	CountryCodeSRB = ISO3166CountryCode{
		Alpha2:      "RS",
		Alpha3:      "SRB",
		Numeric:     "688",
		EnglishName: "Serbia",
		DXCC:        []CountryEnum{CountrySerbia},
	}

	CountryCodeSYC = ISO3166CountryCode{
		Alpha2:      "SC",
		Alpha3:      "SYC",
		Numeric:     "690",
		EnglishName: "Seychelles",
		DXCC:        []CountryEnum{CountrySeychelles},
	}

	CountryCodeSLE = ISO3166CountryCode{
		Alpha2:      "SL",
		Alpha3:      "SLE",
		Numeric:     "694",
		EnglishName: "Sierra Leone",
		DXCC:        []CountryEnum{CountrySierraLeone},
	}

	CountryCodeSGP = ISO3166CountryCode{
		Alpha2:      "SG",
		Alpha3:      "SGP",
		Numeric:     "702",
		EnglishName: "Singapore",
		DXCC:        []CountryEnum{CountrySingapore},
	}

	CountryCodeSXM = ISO3166CountryCode{
		Alpha2:      "SX",
		Alpha3:      "SXM",
		Numeric:     "534",
		EnglishName: "Sint Maarten (Dutch part)",
		DXCC:        []CountryEnum{CountrySintMaarten},
	}

	CountryCodeSVK = ISO3166CountryCode{
		Alpha2:      "SK",
		Alpha3:      "SVK",
		Numeric:     "703",
		EnglishName: "Slovakia",
		DXCC:        []CountryEnum{CountrySlovakRepublic},
	}

	CountryCodeSVN = ISO3166CountryCode{
		Alpha2:      "SI",
		Alpha3:      "SVN",
		Numeric:     "705",
		EnglishName: "Slovenia",
		DXCC:        []CountryEnum{CountrySlovenia},
	}

	CountryCodeSLB = ISO3166CountryCode{
		Alpha2:       "SB",
		Alpha3:       "SLB",
		Numeric:      "090",
		EnglishName:  "Solomon Islands",
		DXCC:         []CountryEnum{CountrySolomonIslands, CountryTemotuProvince},
		Subdivisions: map[string]CountryEnum{"TE": CountryTemotuProvince},
	}

	CountryCodeSOM = ISO3166CountryCode{
		Alpha2:      "SO",
		Alpha3:      "SOM",
		Numeric:     "706",
		EnglishName: "Somalia",
		DXCC:        []CountryEnum{CountrySomalia},
	}

	CountryCodeZAF = ISO3166CountryCode{
		Alpha2:      "ZA",
		Alpha3:      "ZAF",
		Numeric:     "710",
		EnglishName: "South Africa",
		DXCC:        []CountryEnum{CountryRepublicOfSouthAfrica},
	}

	CountryCodeSGS = ISO3166CountryCode{
		Alpha2:      "GS",
		Alpha3:      "SGS",
		Numeric:     "239",
		EnglishName: "South Georgia and the South Sandwich Islands",
		DXCC:        []CountryEnum{CountrySouthGeorgiaIsland, CountrySouthSandwichIslands},
	}

	CountryCodeSSD = ISO3166CountryCode{
		Alpha2:      "SS",
		Alpha3:      "SSD",
		Numeric:     "728",
		EnglishName: "South Sudan",
		DXCC:        []CountryEnum{CountrySouthSudanRepublicOf},
	}

	CountryCodeESP = ISO3166CountryCode{
		Alpha2:      "ES",
		Alpha3:      "ESP",
		Numeric:     "724",
		EnglishName: "Spain",
		DXCC:        []CountryEnum{CountrySpain, CountryBalearicIslands, CountryCeutaMelilla, CountryCanaryIslands},
		Subdivisions: map[string]CountryEnum{
			"IB": CountryBalearicIslands, // Illes Balears autonomous community
			"PM": CountryBalearicIslands, // Illes Balears province
			"CN": CountryCanaryIslands,   // Canarias autonomous community
			"GC": CountryCanaryIslands,   // Las Plamas province
			"TF": CountryCanaryIslands,   // Santa Cruz de Tenerife province
			"CE": CountryCeutaMelilla,    // Ceuta
			"ML": CountryCeutaMelilla,    // Melilla
		},
	}

	CountryCodeLKA = ISO3166CountryCode{
		Alpha2:      "LK",
		Alpha3:      "LKA",
		Numeric:     "144",
		EnglishName: "Sri Lanka",
		DXCC:        []CountryEnum{CountrySriLanka},
	}

	CountryCodeSDN = ISO3166CountryCode{
		Alpha2:      "SD",
		Alpha3:      "SDN",
		Numeric:     "729",
		EnglishName: "Sudan (the)",
		DXCC:        []CountryEnum{CountrySudan},
	}

	CountryCodeSUR = ISO3166CountryCode{
		Alpha2:      "SR",
		Alpha3:      "SUR",
		Numeric:     "740",
		EnglishName: "Suriname",
		DXCC:        []CountryEnum{CountrySuriname},
	}

	CountryCodeSJM = ISO3166CountryCode{
		Alpha2:      "SJ",
		Alpha3:      "SJM",
		Numeric:     "744",
		EnglishName: "Svalbard and Jan Mayen",
		DXCC:        []CountryEnum{CountrySvalbard, CountryJanMayen},
	}

	CountryCodeSWE = ISO3166CountryCode{
		Alpha2:      "SE",
		Alpha3:      "SWE",
		Numeric:     "752",
		EnglishName: "Sweden",
		DXCC:        []CountryEnum{CountrySweden},
	}

	CountryCodeCHE = ISO3166CountryCode{
		Alpha2:      "CH",
		Alpha3:      "CHE",
		Numeric:     "756",
		EnglishName: "Switzerland",
		DXCC:        []CountryEnum{CountrySwitzerland},
	}

	CountryCodeSYR = ISO3166CountryCode{
		Alpha2:      "SY",
		Alpha3:      "SYR",
		Numeric:     "760",
		EnglishName: "Syrian Arab Republic",
		DXCC:        []CountryEnum{CountrySyria},
	}

	CountryCodeTWN = ISO3166CountryCode{
		Alpha2:      "TW",
		Alpha3:      "TWN",
		Numeric:     "158",
		EnglishName: "Taiwan (Province of China)",
		DXCC:        []CountryEnum{CountryTaiwan},
	}

	CountryCodeTJK = ISO3166CountryCode{
		Alpha2:      "TJ",
		Alpha3:      "TJK",
		Numeric:     "762",
		EnglishName: "Tajikistan",
		DXCC:        []CountryEnum{CountryTajikistan},
	}

	CountryCodeTZA = ISO3166CountryCode{
		Alpha2:      "TZ",
		Alpha3:      "TZA",
		Numeric:     "834",
		EnglishName: "Tanzania, United Republic of",
		DXCC:        []CountryEnum{CountryTanzania},
	}

	CountryCodeTHA = ISO3166CountryCode{
		Alpha2:      "TH",
		Alpha3:      "THA",
		Numeric:     "764",
		EnglishName: "Thailand",
		DXCC:        []CountryEnum{CountryThailand},
	}

	CountryCodeTLS = ISO3166CountryCode{
		Alpha2:      "TL",
		Alpha3:      "TLS",
		Numeric:     "626",
		EnglishName: "Timor-Leste",
		DXCC:        []CountryEnum{CountryTimorLeste},
	}

	CountryCodeTGO = ISO3166CountryCode{
		Alpha2:      "TG",
		Alpha3:      "TGO",
		Numeric:     "768",
		EnglishName: "Togo",
		DXCC:        []CountryEnum{CountryTogo},
	}

	CountryCodeTKL = ISO3166CountryCode{
		Alpha2:      "TK",
		Alpha3:      "TKL",
		Numeric:     "772",
		EnglishName: "Tokelau",
		DXCC:        []CountryEnum{CountryTokelauIslands},
	}

	CountryCodeTON = ISO3166CountryCode{
		Alpha2:      "TO",
		Alpha3:      "TON",
		Numeric:     "776",
		EnglishName: "Tonga",
		DXCC:        []CountryEnum{CountryTonga},
	}

	CountryCodeTTO = ISO3166CountryCode{
		Alpha2:      "TT",
		Alpha3:      "TTO",
		Numeric:     "780",
		EnglishName: "Trinidad and Tobago",
		DXCC:        []CountryEnum{CountryTrinidadTobago},
	}

	CountryCodeTUN = ISO3166CountryCode{
		Alpha2:      "TN",
		Alpha3:      "TUN",
		Numeric:     "788",
		EnglishName: "Tunisia",
		DXCC:        []CountryEnum{CountryTunisia},
	}

	CountryCodeTUR = ISO3166CountryCode{
		Alpha2:      "TR",
		Alpha3:      "TUR",
		Numeric:     "792",
		EnglishName: "Turkey",
		DXCC:        []CountryEnum{CountryTurkey},
	}

	CountryCodeTKM = ISO3166CountryCode{
		Alpha2:      "TM",
		Alpha3:      "TKM",
		Numeric:     "795",
		EnglishName: "Turkmenistan",
		DXCC:        []CountryEnum{CountryTurkmenistan},
	}

	CountryCodeTCA = ISO3166CountryCode{
		Alpha2:      "TC",
		Alpha3:      "TCA",
		Numeric:     "796",
		EnglishName: "Turks and Caicos Islands (the)",
		DXCC:        []CountryEnum{CountryTurksCaicosIslands},
	}

	CountryCodeTUV = ISO3166CountryCode{
		Alpha2:      "TV",
		Alpha3:      "TUV",
		Numeric:     "798",
		EnglishName: "Tuvalu",
		DXCC:        []CountryEnum{CountryTuvalu},
	}

	CountryCodeUGA = ISO3166CountryCode{
		Alpha2:      "UG",
		Alpha3:      "UGA",
		Numeric:     "800",
		EnglishName: "Uganda",
		DXCC:        []CountryEnum{CountryUganda},
	}

	CountryCodeUKR = ISO3166CountryCode{
		Alpha2:      "UA",
		Alpha3:      "UKR",
		Numeric:     "804",
		EnglishName: "Ukraine",
		DXCC:        []CountryEnum{CountryUkraine},
	}

	CountryCodeARE = ISO3166CountryCode{
		Alpha2:      "AE",
		Alpha3:      "ARE",
		Numeric:     "784",
		EnglishName: "United Arab Emirates (the)",
		DXCC:        []CountryEnum{CountryUnitedArabEmirates},
	}

	CountryCodeGBR = ISO3166CountryCode{
		Alpha2:      "GB",
		Alpha3:      "GBR",
		Numeric:     "826",
		EnglishName: "United Kingdom of Great Britain and Northern Ireland (the)",
		DXCC:        []CountryEnum{CountryEngland, CountryWales, CountryScotland, CountryNorthernIreland},
		Subdivisions: map[string]CountryEnum{
			// Countries
			"ENG": CountryEngland,
			"NIR": CountryNorthernIreland,
			"SCT": CountryScotland,
			"WLS": CountryWales,
			// English counties, districs, boroughs, authorities
			"BAS": CountryEngland, // Bath and North East Somerset
			"BBD": CountryEngland, // Blackburn with Darwen
			"BCP": CountryEngland, // Bournemouth, Christchurch and Poole
			"BDF": CountryEngland, // Bedford
			"BDG": CountryEngland, // Barking and Dagenham
			"BEN": CountryEngland, // Brent
			"BEX": CountryEngland, // Bexley
			"BIR": CountryEngland, // Birmingham
			"BKM": CountryEngland, // Buckinghamshire
			"BNE": CountryEngland, // Barnet
			"BNH": CountryEngland, // Brighton and Hove
			"BNS": CountryEngland, // Barnsley
			"BOL": CountryEngland, // Bolton
			"BPL": CountryEngland, // Blackpool
			"BRC": CountryEngland, // Bracknell Forest
			"BRD": CountryEngland, // Bradford
			"BRY": CountryEngland, // Bromley
			"BST": CountryEngland, // Bristol, City of
			"BUR": CountryEngland, // Bury
			"CAM": CountryEngland, // Cambridgeshire
			"CBF": CountryEngland, // Central Bedfordshire
			"CHE": CountryEngland, // Cheshire East
			"CHW": CountryEngland, // Cheshire West and Chester
			"CLD": CountryEngland, // Calderdale
			"CMA": CountryEngland, // Cumbria
			"CMD": CountryEngland, // Camden
			"CON": CountryEngland, // Cornwall
			"COV": CountryEngland, // Coventry
			"CRY": CountryEngland, // Croydon
			"DAL": CountryEngland, // Darlington
			"DBY": CountryEngland, // Derbyshire
			"DER": CountryEngland, // Derby
			"DEV": CountryEngland, // Devon
			"DNC": CountryEngland, // Doncaster
			"DOR": CountryEngland, // Dorset
			"DUD": CountryEngland, // Dudley
			"DUR": CountryEngland, // Durham, County
			"EAL": CountryEngland, // Ealing
			"ENF": CountryEngland, // Enfield
			"ERY": CountryEngland, // East Riding of Yorkshire
			"ESS": CountryEngland, // Essex
			"ESX": CountryEngland, // East Sussex
			"GAT": CountryEngland, // Gateshead
			"GLS": CountryEngland, // Gloucestershire
			"GRE": CountryEngland, // Greenwich
			"HAL": CountryEngland, // Halton
			"HAM": CountryEngland, // Hampshire
			"HAV": CountryEngland, // Havering
			"HCK": CountryEngland, // Hackney
			"HEF": CountryEngland, // Herefordshire
			"HIL": CountryEngland, // Hillingdon
			"HMF": CountryEngland, // Hammersmith and Fulham
			"HNS": CountryEngland, // Hounslow
			"HPL": CountryEngland, // Hartlepool
			"HRT": CountryEngland, // Hertfordshire
			"HRW": CountryEngland, // Harrow
			"HRY": CountryEngland, // Haringey
			"IOS": CountryEngland, // Isles of Scilly
			"IOW": CountryEngland, // Isle of Wight
			"ISL": CountryEngland, // Islington
			"KEC": CountryEngland, // Kensington and Chelsea
			"KEN": CountryEngland, // Kent
			"KHL": CountryEngland, // Kingston upon Hull
			"KIR": CountryEngland, // Kirklees
			"KTT": CountryEngland, // Kingston upon Thames
			"KWL": CountryEngland, // Knowsley
			"LAN": CountryEngland, // Lancashire
			"LBH": CountryEngland, // Lambeth
			"LCE": CountryEngland, // Leicester
			"LDS": CountryEngland, // Leeds
			"LEC": CountryEngland, // Leicestershire
			"LEW": CountryEngland, // Lewisham
			"LIN": CountryEngland, // Lincolnshire
			"LIV": CountryEngland, // Liverpool
			"LND": CountryEngland, // London, City of
			"LUT": CountryEngland, // Luton
			"MAN": CountryEngland, // Manchester
			"MDB": CountryEngland, // Middlesbrough
			"MDW": CountryEngland, // Medway
			"MIK": CountryEngland, // Milton Keynes
			"MRT": CountryEngland, // Merton
			"NBL": CountryEngland, // Northumberland
			"NEL": CountryEngland, // North East Lincolnshire
			"NET": CountryEngland, // Newcastle upon Tyne
			"NFK": CountryEngland, // Norfolk
			"NGM": CountryEngland, // Nottingham
			"NLN": CountryEngland, // North Lincolnshire
			"NNH": CountryEngland, // North Northamptonshire
			"NSM": CountryEngland, // North Somerset
			"NTT": CountryEngland, // Nottinghamshire
			"NTY": CountryEngland, // North Tyneside
			"NWM": CountryEngland, // Newham
			"NYK": CountryEngland, // North Yorkshire
			"OLD": CountryEngland, // Oldham
			"OXF": CountryEngland, // Oxfordshire
			"PLY": CountryEngland, // Plymouth
			"POR": CountryEngland, // Portsmouth
			"PTE": CountryEngland, // Peterborough
			"RCC": CountryEngland, // Redcar and Cleveland
			"RCH": CountryEngland, // Rochdale
			"RDB": CountryEngland, // Redbridge
			"RDG": CountryEngland, // Reading
			"RIC": CountryEngland, // Richmond upon Thames
			"ROT": CountryEngland, // Rotherham
			"RUT": CountryEngland, // Rutland
			"SAW": CountryEngland, // Sandwell
			"SFK": CountryEngland, // Suffolk
			"SFT": CountryEngland, // Sefton
			"SGC": CountryEngland, // South Gloucestershire
			"SHF": CountryEngland, // Sheffield
			"SHN": CountryEngland, // St. Helens
			"SHR": CountryEngland, // Shropshire
			"SKP": CountryEngland, // Stockport
			"SLF": CountryEngland, // Salford
			"SLG": CountryEngland, // Slough
			"SND": CountryEngland, // Sunderland
			"SOL": CountryEngland, // Solihull
			"SOM": CountryEngland, // Somerset
			"SOS": CountryEngland, // Southend-on-Sea
			"SRY": CountryEngland, // Surrey
			"STE": CountryEngland, // Stoke-on-Trent
			"STH": CountryEngland, // Southampton
			"STN": CountryEngland, // Sutton
			"STS": CountryEngland, // Staffordshire
			"STT": CountryEngland, // Stockton-on-Tees
			"STY": CountryEngland, // South Tyneside
			"SWD": CountryEngland, // Swindon
			"SWK": CountryEngland, // Southwark
			"TAM": CountryEngland, // Tameside
			"TFW": CountryEngland, // Telford and Wrekin
			"THR": CountryEngland, // Thurrock
			"TOB": CountryEngland, // Torbay
			"TRF": CountryEngland, // Trafford
			"TWH": CountryEngland, // Tower Hamlets
			"WAR": CountryEngland, // Warwickshire
			"WBK": CountryEngland, // West Berkshire
			"WFT": CountryEngland, // Waltham Forest
			"WGN": CountryEngland, // Wigan
			"WIL": CountryEngland, // Wiltshire
			"WKF": CountryEngland, // Wakefield
			"WLL": CountryEngland, // Walsall
			"WLV": CountryEngland, // Wolverhampton
			"WND": CountryEngland, // Wandsworth
			"WNH": CountryEngland, // West Northamptonshire
			"WNM": CountryEngland, // Windsor and Maidenhead
			"WOK": CountryEngland, // Wokingham
			"WOR": CountryEngland, // Worcestershire
			"WRL": CountryEngland, // Wirral
			"WRT": CountryEngland, // Warrington
			"WSM": CountryEngland, // Westminster
			"WSX": CountryEngland, // West Sussex
			"YOR": CountryEngland, // York
			// Northern Irish districs
			"ABC": CountryNorthernIreland, // Armagh, Banbridge, and Craigavon
			"AND": CountryNorthernIreland, // Ards and North Down
			"ANN": CountryNorthernIreland, // Antrim and Newtonabbey
			"BFS": CountryNorthernIreland, // Belfast City
			"CCG": CountryNorthernIreland, // Causeway Coast and Glens
			"DRS": CountryNorthernIreland, // Derry and Strabane
			"FMO": CountryNorthernIreland, // Fermanagh and Omagh
			"LBC": CountryNorthernIreland, // Lisburn and Castlereagh
			"MEA": CountryNorthernIreland, // Mid and East Antrim
			"MUL": CountryNorthernIreland, // Mid-Ulster
			"NMD": CountryNorthernIreland, // Newry, Mourne, and Down
			// Scottish council areas
			"ABD": CountryScotland, // Aberdeenshire
			"ABE": CountryScotland, // Aberdeen City
			"AGB": CountryScotland, // Argyll and Bute
			"ANS": CountryScotland, // Angus
			"CLK": CountryScotland, // Clackmannanshire
			"DGY": CountryScotland, // Dumfries and Galloway
			"DND": CountryScotland, // Dundee City
			"EAY": CountryScotland, // East Ayrshire
			"EDH": CountryScotland, // Edinburgh, City of
			"EDU": CountryScotland, // East Dunbartonshire
			"ELN": CountryScotland, // East Lothian
			"ELS": CountryScotland, // Eilean Siar
			"ERW": CountryScotland, // East Renfrewshire
			"FAL": CountryScotland, // Falkirk
			"FIF": CountryScotland, // Fife
			"GLG": CountryScotland, // Glasgow City
			"HLD": CountryScotland, // Highland
			"IVC": CountryScotland, // Inverclyde
			"MLN": CountryScotland, // Midlothian
			"MRY": CountryScotland, // Moray
			"NAY": CountryScotland, // North Ayrshire
			"NLK": CountryScotland, // North Lanarkshire
			"ORK": CountryScotland, // Orkney Islands
			"PKN": CountryScotland, // Perth and Kinross
			"RFW": CountryScotland, // Renfrewshire
			"SAY": CountryScotland, // South Ayrshire
			"SCB": CountryScotland, // Scottish Borders
			"SLK": CountryScotland, // South Lanarkshire
			"STG": CountryScotland, // Stirling
			"WDU": CountryScotland, // West Dunbartonshire
			"WLN": CountryScotland, // West Lothian
			"ZET": CountryScotland, // Shetland Islands
			// Welsh authorities
			"AGY": CountryWales, // Isle of Anglesey
			"BGE": CountryWales, // Bridgend
			"BGW": CountryWales, // Blaenau Gwent
			"CAY": CountryWales, // Caerphilly
			"CGN": CountryWales, // Ceredigion
			"CMN": CountryWales, // Carmarthenshire
			"CRF": CountryWales, // Cardiff
			"CWY": CountryWales, // Conwy
			"DEN": CountryWales, // Denbighshire
			"FLN": CountryWales, // Flintshire
			"GWN": CountryWales, // Gwynedd
			"MON": CountryWales, // Monmouthshire
			"MTY": CountryWales, // Merthyr Tydfil
			"NTL": CountryWales, // Neath Port Talbot
			"NWP": CountryWales, // Newport
			"PEM": CountryWales, // Pembrokeshire
			"POW": CountryWales, // Powys
			"RCT": CountryWales, // Rhondda Cynon Taff
			"SWA": CountryWales, // Swansea
			"TOF": CountryWales, // Torfaen
			"VGL": CountryWales, // Vale of Glamorgan, The
			"WRX": CountryWales, // Wrexham
			// Welsh-language versions of Welsh authorities
			"YNM": CountryWales, // Sir Ynys Môn, AGY
			"POG": CountryWales, // Pen-y-bont ar Ogwr, BGE
			"CAF": CountryWales, // Caerffili, CAY
			"GFY": CountryWales, // Sir Gaerfyrddin, CMN
			"CRD": CountryWales, // Caerdydd, CRF
			"DDB": CountryWales, // Sir Ddinbych, DEN
			"FFL": CountryWales, // Sir y Fflint, FLN
			"FYN": CountryWales, // Sir Fynwy, MON
			"MTU": CountryWales, // Merthyr Tudful, MTY
			"CTL": CountryWales, // Castell-nedd Port Talbot, NTL
			"CNW": CountryWales, // Casnewydd, NWP
			"BNF": CountryWales, // Sir Benfro, PEM
			"ATA": CountryWales, // Abertawe, SWA
			"BMG": CountryWales, // Bro Morgannwg, VGL
			"WRC": CountryWales, // Wrecsam, WRX
		},
	}

	CountryCodeUMI = ISO3166CountryCode{
		Alpha2:      "UM",
		Alpha3:      "UMI",
		Numeric:     "581",
		EnglishName: "United States Minor Outlying Islands (the)",
		DXCC:        []CountryEnum{CountryMidwayIsland, CountryJohnstonIsland, CountryPalmyraJarvisIslands, CountryBakerHowlandIslands, CountryWakeIsland, CountryNavassaIsland},
	}

	CountryCodeUSA = ISO3166CountryCode{
		Alpha2:      "US",
		Alpha3:      "USA",
		Numeric:     "840",
		EnglishName: "United States of America (the)",
		DXCC:        []CountryEnum{CountryUnitedStatesOfAmerica, CountryAlaska, CountryHawaii},
		Subdivisions: map[string]CountryEnum{
			"AK": CountryAlaska,
			"HI": CountryHawaii,
		},
	}

	CountryCodeURY = ISO3166CountryCode{
		Alpha2:      "UY",
		Alpha3:      "URY",
		Numeric:     "858",
		EnglishName: "Uruguay",
		DXCC:        []CountryEnum{CountryUruguay},
	}

	CountryCodeUZB = ISO3166CountryCode{
		Alpha2:      "UZ",
		Alpha3:      "UZB",
		Numeric:     "860",
		EnglishName: "Uzbekistan",
		DXCC:        []CountryEnum{CountryUzbekistan},
	}

	CountryCodeVUT = ISO3166CountryCode{
		Alpha2:      "VU",
		Alpha3:      "VUT",
		Numeric:     "548",
		EnglishName: "Vanuatu",
		DXCC:        []CountryEnum{CountryVanuatu},
	}

	CountryCodeVEN = ISO3166CountryCode{
		Alpha2:      "VE",
		Alpha3:      "VEN",
		Numeric:     "862",
		EnglishName: "Venezuela (Bolivarian Republic of)",
		DXCC:        []CountryEnum{CountryVenezuela},
	}

	CountryCodeVNM = ISO3166CountryCode{
		Alpha2:      "VN",
		Alpha3:      "VNM",
		Numeric:     "704",
		EnglishName: "Viet Nam",
		DXCC:        []CountryEnum{CountryVietNam},
	}

	CountryCodeVGB = ISO3166CountryCode{
		Alpha2:      "VG",
		Alpha3:      "VGB",
		Numeric:     "092",
		EnglishName: "Virgin Islands (British)",
		DXCC:        []CountryEnum{CountryBritishVirginIslands},
	}

	CountryCodeVIR = ISO3166CountryCode{
		Alpha2:      "VI",
		Alpha3:      "VIR",
		Numeric:     "850",
		EnglishName: "Virgin Islands (U.S.)",
		DXCC:        []CountryEnum{CountryVirginIslands},
	}

	CountryCodeWLF = ISO3166CountryCode{
		Alpha2:      "WF",
		Alpha3:      "WLF",
		Numeric:     "876",
		EnglishName: "Wallis and Futuna",
		DXCC:        []CountryEnum{CountryWallisFutunaIslands},
	}

	CountryCodeESH = ISO3166CountryCode{
		Alpha2:      "EH",
		Alpha3:      "ESH",
		Numeric:     "732",
		EnglishName: "Western Sahara",
		DXCC:        []CountryEnum{CountryWesternSahara},
	}

	CountryCodeYEM = ISO3166CountryCode{
		Alpha2:      "YE",
		Alpha3:      "YEM",
		Numeric:     "887",
		EnglishName: "Yemen",
		DXCC:        []CountryEnum{CountryYemen},
	}

	CountryCodeZMB = ISO3166CountryCode{
		Alpha2:      "ZM",
		Alpha3:      "ZMB",
		Numeric:     "894",
		EnglishName: "Zambia",
		DXCC:        []CountryEnum{CountryZambia},
	}

	CountryCodeZWE = ISO3166CountryCode{
		Alpha2:      "ZW",
		Alpha3:      "ZWE",
		Numeric:     "716",
		EnglishName: "Zimbabwe",
		DXCC:        []CountryEnum{CountryZimbabwe},
	}

	CountryCodeALA = ISO3166CountryCode{
		Alpha2:      "AX",
		Alpha3:      "ALA",
		Numeric:     "248",
		EnglishName: "Åland Islands",
		DXCC:        []CountryEnum{CountryAlandIslands},
	}

	ISO3166Countries = []ISO3166CountryCode{
		CountryCodeAFG,
		CountryCodeALB,
		CountryCodeDZA,
		CountryCodeASM,
		CountryCodeAND,
		CountryCodeAGO,
		CountryCodeAIA,
		CountryCodeATA,
		CountryCodeATG,
		CountryCodeARG,
		CountryCodeARM,
		CountryCodeABW,
		CountryCodeAUS,
		CountryCodeAUT,
		CountryCodeAZE,
		CountryCodeBHS,
		CountryCodeBHR,
		CountryCodeBGD,
		CountryCodeBRB,
		CountryCodeBLR,
		CountryCodeBEL,
		CountryCodeBLZ,
		CountryCodeBEN,
		CountryCodeBMU,
		CountryCodeBTN,
		CountryCodeBOL,
		CountryCodeBES,
		CountryCodeBIH,
		CountryCodeBWA,
		CountryCodeBVT,
		CountryCodeBRA,
		CountryCodeIOT,
		CountryCodeBRN,
		CountryCodeBGR,
		CountryCodeBFA,
		CountryCodeBDI,
		CountryCodeCPV,
		CountryCodeKHM,
		CountryCodeCMR,
		CountryCodeCAN,
		CountryCodeCYM,
		CountryCodeCAF,
		CountryCodeTCD,
		CountryCodeCHL,
		CountryCodeCHN,
		CountryCodeCXR,
		CountryCodeCCK,
		CountryCodeCOL,
		CountryCodeCOM,
		CountryCodeCOD,
		CountryCodeCOG,
		CountryCodeCOK,
		CountryCodeCRI,
		CountryCodeHRV,
		CountryCodeCUB,
		CountryCodeCUW,
		CountryCodeCYP,
		CountryCodeCZE,
		CountryCodeCIV,
		CountryCodeDNK,
		CountryCodeDJI,
		CountryCodeDMA,
		CountryCodeDOM,
		CountryCodeECU,
		CountryCodeEGY,
		CountryCodeSLV,
		CountryCodeGNQ,
		CountryCodeERI,
		CountryCodeEST,
		CountryCodeSWZ,
		CountryCodeETH,
		CountryCodeFLK,
		CountryCodeFRO,
		CountryCodeFJI,
		CountryCodeFIN,
		CountryCodeFRA,
		CountryCodeGUF,
		CountryCodePYF,
		CountryCodeATF,
		CountryCodeGAB,
		CountryCodeGMB,
		CountryCodeGEO,
		CountryCodeDEU,
		CountryCodeGHA,
		CountryCodeGIB,
		CountryCodeGRC,
		CountryCodeGRL,
		CountryCodeGRD,
		CountryCodeGLP,
		CountryCodeGUM,
		CountryCodeGTM,
		CountryCodeGGY,
		CountryCodeGIN,
		CountryCodeGNB,
		CountryCodeGUY,
		CountryCodeHTI,
		CountryCodeHMD,
		CountryCodeVAT,
		CountryCodeHND,
		CountryCodeHKG,
		CountryCodeHUN,
		CountryCodeISL,
		CountryCodeIND,
		CountryCodeIDN,
		CountryCodeIRN,
		CountryCodeIRQ,
		CountryCodeIRL,
		CountryCodeIMN,
		CountryCodeISR,
		CountryCodeITA,
		CountryCodeJAM,
		CountryCodeJPN,
		CountryCodeJEY,
		CountryCodeJOR,
		CountryCodeKAZ,
		CountryCodeKEN,
		CountryCodeKIR,
		CountryCodePRK,
		CountryCodeKOR,
		CountryCodeKWT,
		CountryCodeKGZ,
		CountryCodeLAO,
		CountryCodeLVA,
		CountryCodeLBN,
		CountryCodeLSO,
		CountryCodeLBR,
		CountryCodeLBY,
		CountryCodeLIE,
		CountryCodeLTU,
		CountryCodeLUX,
		CountryCodeMAC,
		CountryCodeMDG,
		CountryCodeMWI,
		CountryCodeMYS,
		CountryCodeMDV,
		CountryCodeMLI,
		CountryCodeMLT,
		CountryCodeMHL,
		CountryCodeMTQ,
		CountryCodeMRT,
		CountryCodeMUS,
		CountryCodeMYT,
		CountryCodeMEX,
		CountryCodeFSM,
		CountryCodeMDA,
		CountryCodeMCO,
		CountryCodeMNG,
		CountryCodeMNE,
		CountryCodeMSR,
		CountryCodeMAR,
		CountryCodeMOZ,
		CountryCodeMMR,
		CountryCodeNAM,
		CountryCodeNRU,
		CountryCodeNPL,
		CountryCodeNLD,
		CountryCodeNCL,
		CountryCodeNZL,
		CountryCodeNIC,
		CountryCodeNER,
		CountryCodeNGA,
		CountryCodeNIU,
		CountryCodeNFK,
		CountryCodeMNP,
		CountryCodeNOR,
		CountryCodeOMN,
		CountryCodePAK,
		CountryCodePLW,
		CountryCodePSE,
		CountryCodePAN,
		CountryCodePNG,
		CountryCodePRY,
		CountryCodePER,
		CountryCodePHL,
		CountryCodePCN,
		CountryCodePOL,
		CountryCodePRT,
		CountryCodePRI,
		CountryCodeQAT,
		CountryCodeMKD,
		CountryCodeROU,
		CountryCodeRUS,
		CountryCodeRWA,
		CountryCodeREU,
		CountryCodeBLM,
		CountryCodeSHN,
		CountryCodeKNA,
		CountryCodeLCA,
		CountryCodeMAF,
		CountryCodeSPM,
		CountryCodeVCT,
		CountryCodeWSM,
		CountryCodeSMR,
		CountryCodeSTP,
		CountryCodeSAU,
		CountryCodeSEN,
		CountryCodeSRB,
		CountryCodeSYC,
		CountryCodeSLE,
		CountryCodeSGP,
		CountryCodeSXM,
		CountryCodeSVK,
		CountryCodeSVN,
		CountryCodeSLB,
		CountryCodeSOM,
		CountryCodeZAF,
		CountryCodeSGS,
		CountryCodeSSD,
		CountryCodeESP,
		CountryCodeLKA,
		CountryCodeSDN,
		CountryCodeSUR,
		CountryCodeSJM,
		CountryCodeSWE,
		CountryCodeCHE,
		CountryCodeSYR,
		CountryCodeTWN,
		CountryCodeTJK,
		CountryCodeTZA,
		CountryCodeTHA,
		CountryCodeTLS,
		CountryCodeTGO,
		CountryCodeTKL,
		CountryCodeTON,
		CountryCodeTTO,
		CountryCodeTUN,
		CountryCodeTUR,
		CountryCodeTKM,
		CountryCodeTCA,
		CountryCodeTUV,
		CountryCodeUGA,
		CountryCodeUKR,
		CountryCodeARE,
		CountryCodeGBR,
		CountryCodeUMI,
		CountryCodeUSA,
		CountryCodeURY,
		CountryCodeUZB,
		CountryCodeVUT,
		CountryCodeVEN,
		CountryCodeVNM,
		CountryCodeVGB,
		CountryCodeVIR,
		CountryCodeWLF,
		CountryCodeESH,
		CountryCodeYEM,
		CountryCodeZMB,
		CountryCodeZWE,
		CountryCodeALA,
		CountryCodeXKX, // TODO remove when Kosovo gets an assigned ISO code
	}
)
