// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spec

import "strings"

// ITUZoneFor returns a list of ITU Zone numbers for a DXCC entity code or
// country name.  Returns an empty slice if the entity code or country is not
// known.  For some countries which span multiple CQ Zones, the CqZone property
// of the SecondaryAdministrativeSubdivision enumeration indicates the station's
// zone.  This list was compiled from
// https://www.arrl.org/files/file/DXCC/2022_Current_Deleted.txt although it may
// contain geographic errors, e.g. Svalbard should extend into Zone 75.
// See also https://www.mapability.com/ei8ic/maps/ituzone.php and
// https://zone-info.eu/
func ITUZoneFor(name string) []int {
	switch strings.ToUpper(name) {
	default:
		return []int{}

	// Multi-zone countries
	case CountryAlaska.EntityName, CountryAlaska.EntityCode:
		return []int{1, 2}
	case CountryCanada.EntityName, CountryCanada.EntityCode:
		return []int{2, 3, 4, 9, 75}
	case CountryGreenland.EntityName, CountryGreenland.EntityCode:
		return []int{5, 75}
	case CountryUnitedStatesOfAmerica.EntityName, CountryUnitedStatesOfAmerica.EntityCode:
		return []int{6, 7, 8}
	case CountryBrazil.EntityName, CountryBrazil.EntityCode:
		return []int{12, 13, 15}
	case CountryBolivia.EntityName, CountryBolivia.EntityCode:
		return []int{12, 14}
	case CountryArgentina.EntityName, CountryArgentina.EntityCode:
		return []int{14, 16}
	case CountryChile.EntityName, CountryChile.EntityCode:
		return []int{14, 16}
	case CountryEuropeanRussia.EntityName, CountryEuropeanRussia.EntityCode:
		return []int{19, 20, 29, 30}
	case CountryAsiaticRussia.EntityName, CountryAsiaticRussia.EntityCode:
		return []int{20, 21, 22, 23, 24, 25, 26, 30, 31, 32, 33, 34, 35, 75}
	case CountryKazakhstan.EntityName, CountryKazakhstan.EntityCode:
		return []int{29, 30, 31}
	case CountryKyrgyzstan.EntityName, CountryKyrgyzstan.EntityCode:
		return []int{30, 31}
	case CountryMongolia.EntityName, CountryMongolia.EntityCode:
		return []int{32, 33}
	case CountryChina.EntityName, CountryChina.EntityCode:
		return []int{33, 42, 43, 44}
	case CountrySouthernSudan_deleted.EntityName, CountrySouthernSudan_deleted.EntityCode:
		return []int{47, 48}
	case CountrySudan.EntityName, CountrySudan.EntityCode:
		return []int{47, 48}
	case CountryFrenchEquatorialAfrica_deleted.EntityName, CountryFrenchEquatorialAfrica_deleted.EntityCode:
		return []int{47, 52}
	case CountryIndonesia.EntityName, CountryIndonesia.EntityCode:
		return []int{51, 54}
	case CountryAustralia.EntityName, CountryAustralia.EntityCode:
		return []int{55, 58, 59}
	case CountryPalmyraJarvisIslands.EntityName, CountryPalmyraJarvisIslands.EntityCode:
		return []int{61, 62}
	case CountryEKiribatiLineIslands.EntityName, CountryEKiribatiLineIslands.EntityCode:
		return []int{61, 63}
	case CountryAntarctica.EntityName, CountryAntarctica.EntityCode:
		return []int{67, 69, 70, 71, 72, 73, 74}

	// Zones 1 through 7 are just Alaska, Canada, and USA
	// Zone 8
	case CountryUnitedNationsHq.EntityName, CountryUnitedNationsHq.EntityCode:
		return []int{8}

	// Zone 9
	case CountryNewfoundlandLabrador_deleted.EntityName, CountryNewfoundlandLabrador_deleted.EntityCode,
		CountrySableIsland.EntityName, CountrySableIsland.EntityCode,
		CountryStPaulIsland.EntityName, CountryStPaulIsland.EntityCode,
		CountryStPierreMiquelon.EntityName, CountryStPierreMiquelon.EntityCode:
		return []int{9}

	// Zone 10
	case CountryClippertonIsland.EntityName, CountryClippertonIsland.EntityCode,
		CountryMexico.EntityName, CountryMexico.EntityCode,
		CountryRevillagigedo.EntityName, CountryRevillagigedo.EntityCode:
		return []int{10}

	// Zone 11
	case CountryAnguilla.EntityName, CountryAnguilla.EntityCode,
		CountryAntiguaBarbuda.EntityName, CountryAntiguaBarbuda.EntityCode,
		CountryAruba.EntityName, CountryAruba.EntityCode,
		CountryAvesIsland.EntityName, CountryAvesIsland.EntityCode,
		CountryBahamas.EntityName, CountryBahamas.EntityCode,
		CountryBajoNuevo_deleted.EntityName, CountryBajoNuevo_deleted.EntityCode,
		CountryBarbados.EntityName, CountryBarbados.EntityCode,
		CountryBelize.EntityName, CountryBelize.EntityCode,
		CountryBermuda.EntityName, CountryBermuda.EntityCode,
		CountryBonaire.EntityName, CountryBonaire.EntityCode,
		CountryBonaireCuracao_deleted.EntityName, CountryBonaireCuracao_deleted.EntityCode,
		CountryBritishVirginIslands.EntityName, CountryBritishVirginIslands.EntityCode,
		CountryCanalZone_deleted.EntityName, CountryCanalZone_deleted.EntityCode,
		CountryCaymanIslands.EntityName, CountryCaymanIslands.EntityCode,
		CountryCostaRica.EntityName, CountryCostaRica.EntityCode,
		CountryCuba.EntityName, CountryCuba.EntityCode,
		CountryCuracao.EntityName, CountryCuracao.EntityCode,
		CountryDesecheoIsland.EntityName, CountryDesecheoIsland.EntityCode,
		CountryDominica.EntityName, CountryDominica.EntityCode,
		CountryDominicanRepublic.EntityName, CountryDominicanRepublic.EntityCode,
		CountryElSalvador.EntityName, CountryElSalvador.EntityCode,
		CountryGrenada.EntityName, CountryGrenada.EntityCode,
		CountryGuadeloupe.EntityName, CountryGuadeloupe.EntityCode,
		CountryGuantanamoBay.EntityName, CountryGuantanamoBay.EntityCode,
		CountryHaiti.EntityName, CountryHaiti.EntityCode,
		CountryHonduras.EntityName, CountryHonduras.EntityCode,
		CountryJamaica.EntityName, CountryJamaica.EntityCode,
		CountryMartinique.EntityName, CountryMartinique.EntityCode,
		CountryMontserrat.EntityName, CountryMontserrat.EntityCode,
		CountryNavassaIsland.EntityName, CountryNavassaIsland.EntityCode,
		CountryNicaragua.EntityName, CountryNicaragua.EntityCode,
		CountryPanama.EntityName, CountryPanama.EntityCode,
		CountryPuertoRico.EntityName, CountryPuertoRico.EntityCode,
		CountrySabaStEustatius.EntityName, CountrySabaStEustatius.EntityCode,
		CountrySaintBarthelemy.EntityName, CountrySaintBarthelemy.EntityCode,
		CountrySaintMartin.EntityName, CountrySaintMartin.EntityCode,
		CountrySanAndresProvidencia.EntityName, CountrySanAndresProvidencia.EntityCode,
		CountrySerranaBankRoncadorCay_deleted.EntityName, CountrySerranaBankRoncadorCay_deleted.EntityCode,
		CountrySintMaarten.EntityName, CountrySintMaarten.EntityCode,
		CountryStKittsNevis.EntityName, CountryStKittsNevis.EntityCode,
		CountryStLucia.EntityName, CountryStLucia.EntityCode,
		CountryStMaartenSabaStEustatius_deleted.EntityName, CountryStMaartenSabaStEustatius_deleted.EntityCode,
		CountryStVincent.EntityName, CountryStVincent.EntityCode,
		CountrySwanIslands_deleted.EntityName, CountrySwanIslands_deleted.EntityCode,
		CountryTrinidadTobago.EntityName, CountryTrinidadTobago.EntityCode,
		CountryTurksCaicosIslands.EntityName, CountryTurksCaicosIslands.EntityCode,
		CountryVirginIslands.EntityName, CountryVirginIslands.EntityCode:
		return []int{11}

	// Zone 12
	case CountryCocosIsland.EntityName, CountryCocosIsland.EntityCode,
		CountryColombia.EntityName, CountryColombia.EntityCode,
		CountryEcuador.EntityName, CountryEcuador.EntityCode,
		CountryFrenchGuiana.EntityName, CountryFrenchGuiana.EntityCode,
		CountryGalapagosIslands.EntityName, CountryGalapagosIslands.EntityCode,
		CountryGuatemala.EntityName, CountryGuatemala.EntityCode,
		CountryGuyana.EntityName, CountryGuyana.EntityCode,
		CountryMalpeloIsland.EntityName, CountryMalpeloIsland.EntityCode,
		CountryPeru.EntityName, CountryPeru.EntityCode,
		CountrySuriname.EntityName, CountrySuriname.EntityCode,
		CountryVenezuela.EntityName, CountryVenezuela.EntityCode:
		return []int{12}

	// Zone 13
	case CountryFernandoDeNoronha.EntityName, CountryFernandoDeNoronha.EntityCode:
		return []int{13}
	case CountryStPeterStPaulRocks.EntityName, CountryStPeterStPaulRocks.EntityCode:
		return []int{13}

	// Zone 14
	case CountryJuanFernandezIslands.EntityName, CountryJuanFernandezIslands.EntityCode,
		CountryParaguay.EntityName, CountryParaguay.EntityCode,
		CountrySanFelixSanAmbrosio.EntityName, CountrySanFelixSanAmbrosio.EntityCode,
		CountryUruguay.EntityName, CountryUruguay.EntityCode:
		return []int{14}

	// Zone 15
	case CountryTrindadeMartimVazIslands.EntityName, CountryTrindadeMartimVazIslands.EntityCode:
		return []int{15}

	// Zone 16
	case CountryFalklandIslands.EntityName, CountryFalklandIslands.EntityCode:
		return []int{16}

	// Zone 17
	case CountryIceland.EntityName, CountryIceland.EntityCode:
		return []int{17}

	// Zone 18
	case CountryAlandIslands.EntityName, CountryAlandIslands.EntityCode,
		CountryDenmark.EntityName, CountryDenmark.EntityCode,
		CountryFaroeIslands.EntityName, CountryFaroeIslands.EntityCode,
		CountryFinland.EntityName, CountryFinland.EntityCode,
		CountryJanMayen.EntityName, CountryJanMayen.EntityCode,
		CountryMarketReef.EntityName, CountryMarketReef.EntityCode,
		CountryNorway.EntityName, CountryNorway.EntityCode,
		CountrySvalbard.EntityName, CountrySvalbard.EntityCode,
		CountrySweden.EntityName, CountrySweden.EntityCode:
		return []int{18}

	// Zone 19
	case CountryKareloFinnishRepublic_deleted.EntityName, CountryKareloFinnishRepublic_deleted.EntityCode:
		return []int{19}

	// Zone 20 to 26 are all just Russia
	case CountryAndorra.EntityName, CountryAndorra.EntityCode,
		CountryBelgium.EntityName, CountryBelgium.EntityCode,
		CountryEngland.EntityName, CountryEngland.EntityCode,
		CountryFrance.EntityName, CountryFrance.EntityCode,
		CountryGuernsey.EntityName, CountryGuernsey.EntityCode,
		CountryIreland.EntityName, CountryIreland.EntityCode,
		CountryIsleOfMan.EntityName, CountryIsleOfMan.EntityCode,
		CountryJersey.EntityName, CountryJersey.EntityCode,
		CountryLuxembourg.EntityName, CountryLuxembourg.EntityCode,
		CountryMonaco.EntityName, CountryMonaco.EntityCode,
		CountryNetherlands.EntityName, CountryNetherlands.EntityCode,
		CountryNorthernIreland.EntityName, CountryNorthernIreland.EntityCode,
		CountryScotland.EntityName, CountryScotland.EntityCode,
		CountryWales.EntityName, CountryWales.EntityCode:
		return []int{27}

	// Zone 28
	case CountryAlbania.EntityName, CountryAlbania.EntityCode,
		CountryAustria.EntityName, CountryAustria.EntityCode,
		CountryBosniaHerzegovina.EntityName, CountryBosniaHerzegovina.EntityCode,
		CountryBulgaria.EntityName, CountryBulgaria.EntityCode,
		CountryCorsica.EntityName, CountryCorsica.EntityCode,
		CountryCrete.EntityName, CountryCrete.EntityCode,
		CountryCroatia.EntityName, CountryCroatia.EntityCode,
		CountryCzechRepublic.EntityName, CountryCzechRepublic.EntityCode,
		CountryCzechoslovakia_deleted.EntityName, CountryCzechoslovakia_deleted.EntityCode,
		CountryDodecanese.EntityName, CountryDodecanese.EntityCode,
		CountryFederalRepublicOfGermany.EntityName, CountryFederalRepublicOfGermany.EntityCode,
		CountryGermanDemocraticRepublic_deleted.EntityName, CountryGermanDemocraticRepublic_deleted.EntityCode,
		CountryGermany_deleted.EntityName, CountryGermany_deleted.EntityCode,
		CountryGreece.EntityName, CountryGreece.EntityCode,
		CountryHungary.EntityName, CountryHungary.EntityCode,
		CountryItuHq.EntityName, CountryItuHq.EntityCode,
		CountryItaly.EntityName, CountryItaly.EntityCode,
		CountryLiechtenstein.EntityName, CountryLiechtenstein.EntityCode,
		CountryMalta.EntityName, CountryMalta.EntityCode,
		CountryMontenegro.EntityName, CountryMontenegro.EntityCode,
		CountryMountAthos.EntityName, CountryMountAthos.EntityCode,
		CountryNorthMacedoniaRepublicOf.EntityName, CountryNorthMacedoniaRepublicOf.EntityCode,
		CountryPoland.EntityName, CountryPoland.EntityCode,
		CountryRepublicOfKosovo.EntityName, CountryRepublicOfKosovo.EntityCode,
		CountryRomania.EntityName, CountryRomania.EntityCode,
		CountrySaar_deleted.EntityName, CountrySaar_deleted.EntityCode,
		CountrySanMarino.EntityName, CountrySanMarino.EntityCode,
		CountrySardinia.EntityName, CountrySardinia.EntityCode,
		CountrySerbia.EntityName, CountrySerbia.EntityCode,
		CountrySlovakRepublic.EntityName, CountrySlovakRepublic.EntityCode,
		CountrySlovenia.EntityName, CountrySlovenia.EntityCode,
		CountrySovereignMilitaryOrderOfMalta.EntityName, CountrySovereignMilitaryOrderOfMalta.EntityCode,
		CountrySwitzerland.EntityName, CountrySwitzerland.EntityCode,
		CountryTrieste_deleted.EntityName, CountryTrieste_deleted.EntityCode,
		CountryVatican.EntityName, CountryVatican.EntityCode:
		return []int{28}

	// Zone 29
	case CountryArmenia.EntityName, CountryArmenia.EntityCode,
		CountryAzerbaijan.EntityName, CountryAzerbaijan.EntityCode,
		CountryBelarus.EntityName, CountryBelarus.EntityCode,
		CountryEstonia.EntityName, CountryEstonia.EntityCode,
		CountryGeorgia.EntityName, CountryGeorgia.EntityCode,
		CountryKaliningrad.EntityName, CountryKaliningrad.EntityCode,
		CountryLatvia.EntityName, CountryLatvia.EntityCode,
		CountryLithuania.EntityName, CountryLithuania.EntityCode,
		CountryMalyjVysotskijIsland_deleted.EntityName, CountryMalyjVysotskijIsland_deleted.EntityCode,
		CountryMoldova.EntityName, CountryMoldova.EntityCode,
		CountryUkraine.EntityName, CountryUkraine.EntityCode:
		return []int{29}

	// Zone 30
	case CountryTajikistan.EntityName, CountryTajikistan.EntityCode:
		return []int{30}
	case CountryTurkmenistan.EntityName, CountryTurkmenistan.EntityCode:
		return []int{30}
	case CountryUzbekistan.EntityName, CountryUzbekistan.EntityCode:
		return []int{30}

	// Zones 31 and 32 are all split countries
	// Zone 33
	case CountryManchuria_deleted.EntityName, CountryManchuria_deleted.EntityCode:
		return []int{33}

	// Zones 34 and 35 are just Russia
	// Zone 36
	case CountryAzores.EntityName, CountryAzores.EntityCode,
		CountryCanaryIslands.EntityName, CountryCanaryIslands.EntityCode,
		CountryMadeiraIslands.EntityName, CountryMadeiraIslands.EntityCode:
		return []int{36}

	// Zone 37
	case CountryAlgeria.EntityName, CountryAlgeria.EntityCode,
		CountryBalearicIslands.EntityName, CountryBalearicIslands.EntityCode,
		CountryCeutaMelilla.EntityName, CountryCeutaMelilla.EntityCode,
		CountryGibraltar.EntityName, CountryGibraltar.EntityCode,
		CountryIfni_deleted.EntityName, CountryIfni_deleted.EntityCode,
		CountryMorocco.EntityName, CountryMorocco.EntityCode,
		CountryPortugal.EntityName, CountryPortugal.EntityCode,
		CountrySpain.EntityName, CountrySpain.EntityCode,
		CountryTangier_deleted.EntityName, CountryTangier_deleted.EntityCode,
		CountryTunisia.EntityName, CountryTunisia.EntityCode:
		return []int{37}

	// Zone 38
	case CountryEgypt.EntityName, CountryEgypt.EntityCode,
		CountryLibya.EntityName, CountryLibya.EntityCode:
		return []int{38}

	// Zone 39
	case CountryAbuAilIslands_deleted.EntityName, CountryAbuAilIslands_deleted.EntityCode,
		CountryBahrain.EntityName, CountryBahrain.EntityCode,
		CountryCyprus.EntityName, CountryCyprus.EntityCode,
		CountryIraq.EntityName, CountryIraq.EntityCode,
		CountryIsrael.EntityName, CountryIsrael.EntityCode,
		CountryJordan.EntityName, CountryJordan.EntityCode,
		CountryKamaranIslands_deleted.EntityName, CountryKamaranIslands_deleted.EntityCode,
		CountryKuriaMuriaIsland_deleted.EntityName, CountryKuriaMuriaIsland_deleted.EntityCode,
		CountryKuwait.EntityName, CountryKuwait.EntityCode,
		CountryKuwaitSaudiArabiaNeutralZone_deleted.EntityName, CountryKuwaitSaudiArabiaNeutralZone_deleted.EntityCode,
		CountryLebanon.EntityName, CountryLebanon.EntityCode,
		CountryOman.EntityName, CountryOman.EntityCode,
		CountryPalestine.EntityName, CountryPalestine.EntityCode,
		CountryPalestine_deleted.EntityName, CountryPalestine_deleted.EntityCode,
		CountryPeoplesDemocraticRepOfYemen_deleted.EntityName, CountryPeoplesDemocraticRepOfYemen_deleted.EntityCode,
		CountryQatar.EntityName, CountryQatar.EntityCode,
		CountrySaudiArabia.EntityName, CountrySaudiArabia.EntityCode,
		CountrySaudiArabiaIraqNeutralZone_deleted.EntityName, CountrySaudiArabiaIraqNeutralZone_deleted.EntityCode,
		CountrySyria.EntityName, CountrySyria.EntityCode,
		CountryTurkey.EntityName, CountryTurkey.EntityCode,
		CountryUkSovereignBaseAreasOnCyprus.EntityName, CountryUkSovereignBaseAreasOnCyprus.EntityCode,
		CountryUnitedArabEmirates.EntityName, CountryUnitedArabEmirates.EntityCode,
		CountryYemen.EntityName, CountryYemen.EntityCode,
		CountryYemenArabRepublic_deleted.EntityName, CountryYemenArabRepublic_deleted.EntityCode:
		return []int{39}

	// Zone 40
	case CountryAfghanistan.EntityName, CountryAfghanistan.EntityCode,
		CountryIran.EntityName, CountryIran.EntityCode:
		return []int{40}

	// Zone 41
	case CountryBangladesh.EntityName, CountryBangladesh.EntityCode,
		CountryBhutan.EntityName, CountryBhutan.EntityCode,
		CountryBlenheimReef_deleted.EntityName, CountryBlenheimReef_deleted.EntityCode,
		CountryChagosIslands.EntityName, CountryChagosIslands.EntityCode,
		CountryDamaoDiu_deleted.EntityName, CountryDamaoDiu_deleted.EntityCode,
		CountryFrenchIndia_deleted.EntityName, CountryFrenchIndia_deleted.EntityCode,
		CountryGoa_deleted.EntityName, CountryGoa_deleted.EntityCode,
		CountryIndia.EntityName, CountryIndia.EntityCode,
		CountryLakshadweepIslands.EntityName, CountryLakshadweepIslands.EntityCode,
		CountryMaldives.EntityName, CountryMaldives.EntityCode,
		CountryPakistan.EntityName, CountryPakistan.EntityCode,
		CountrySikkim_deleted.EntityName, CountrySikkim_deleted.EntityCode,
		CountrySriLanka.EntityName, CountrySriLanka.EntityCode,
		CountryTibet_deleted.EntityName, CountryTibet_deleted.EntityCode:
		return []int{41}

	// Zone 42
	case CountryNepal.EntityName, CountryNepal.EntityCode:
		return []int{42}

	// Zone 43 is just China
	// Zone 44
	case CountryDemocraticPeoplesRepOfKorea.EntityName, CountryDemocraticPeoplesRepOfKorea.EntityCode,
		CountryHongKong.EntityName, CountryHongKong.EntityCode,
		CountryMacao.EntityName, CountryMacao.EntityCode,
		CountryPratasIsland.EntityName, CountryPratasIsland.EntityCode,
		CountryRepublicOfKorea.EntityName, CountryRepublicOfKorea.EntityCode,
		CountryTaiwan.EntityName, CountryTaiwan.EntityCode:
		return []int{44}

	// Zone 45
	case CountryJapan.EntityName, CountryJapan.EntityCode,
		CountryOgasawara.EntityName, CountryOgasawara.EntityCode,
		CountryOkinawaRyukyuIslands_deleted.EntityName, CountryOkinawaRyukyuIslands_deleted.EntityCode,
		CountryOkinoToriShima_deleted.EntityName, CountryOkinoToriShima_deleted.EntityCode:
		return []int{45}

	// Zone 46
	case CountryBenin.EntityName, CountryBenin.EntityCode,
		CountryBurkinaFaso.EntityName, CountryBurkinaFaso.EntityCode,
		CountryCapeVerde.EntityName, CountryCapeVerde.EntityCode,
		CountryCoteDIvoire.EntityName, CountryCoteDIvoire.EntityCode,
		CountryFrenchWestAfrica_deleted.EntityName, CountryFrenchWestAfrica_deleted.EntityCode,
		CountryGhana.EntityName, CountryGhana.EntityCode,
		CountryGoldCoastTogoland_deleted.EntityName, CountryGoldCoastTogoland_deleted.EntityCode,
		CountryGuinea.EntityName, CountryGuinea.EntityCode,
		CountryGuineaBissau.EntityName, CountryGuineaBissau.EntityCode,
		CountryLiberia.EntityName, CountryLiberia.EntityCode,
		CountryMali.EntityName, CountryMali.EntityCode,
		CountryMauritania.EntityName, CountryMauritania.EntityCode,
		CountryNiger.EntityName, CountryNiger.EntityCode,
		CountryNigeria.EntityName, CountryNigeria.EntityCode,
		CountrySenegal.EntityName, CountrySenegal.EntityCode,
		CountrySierraLeone.EntityName, CountrySierraLeone.EntityCode,
		CountryTheGambia.EntityName, CountryTheGambia.EntityCode,
		CountryTogo.EntityName, CountryTogo.EntityCode,
		CountryWesternSahara.EntityName, CountryWesternSahara.EntityCode:
		return []int{46}

	// Zone 47
	case CountryCameroon.EntityName, CountryCameroon.EntityCode,
		CountryCentralAfrica.EntityName, CountryCentralAfrica.EntityCode,
		CountryChad.EntityName, CountryChad.EntityCode,
		CountryEquatorialGuinea.EntityName, CountryEquatorialGuinea.EntityCode,
		CountrySaoTomePrincipe.EntityName, CountrySaoTomePrincipe.EntityCode:
		return []int{47}

	// Zone 48
	case CountryBritishSomaliland_deleted.EntityName, CountryBritishSomaliland_deleted.EntityCode,
		CountryDjibouti.EntityName, CountryDjibouti.EntityCode,
		CountryEritrea.EntityName, CountryEritrea.EntityCode,
		CountryEthiopia.EntityName, CountryEthiopia.EntityCode,
		CountryItalianSomaliland_deleted.EntityName, CountryItalianSomaliland_deleted.EntityCode,
		CountryKenya.EntityName, CountryKenya.EntityCode,
		CountrySomalia.EntityName, CountrySomalia.EntityCode,
		CountrySouthSudanRepublicOf.EntityName, CountrySouthSudanRepublicOf.EntityCode,
		CountryUganda.EntityName, CountryUganda.EntityCode:
		return []int{48}

	// Zone 49
	case CountryAndamanNicobarIslands.EntityName, CountryAndamanNicobarIslands.EntityCode,
		CountryCambodia.EntityName, CountryCambodia.EntityCode,
		CountryFrenchIndoChina_deleted.EntityName, CountryFrenchIndoChina_deleted.EntityCode,
		CountryLaos.EntityName, CountryLaos.EntityCode,
		CountryMyanmar.EntityName, CountryMyanmar.EntityCode,
		CountryThailand.EntityName, CountryThailand.EntityCode,
		CountryVietNam.EntityName, CountryVietNam.EntityCode:
		return []int{49}

	// Zone 50
	case CountryPhilippines.EntityName, CountryPhilippines.EntityCode,
		CountryScarboroughReef.EntityName, CountryScarboroughReef.EntityCode,
		CountrySpratlyIslands.EntityName, CountrySpratlyIslands.EntityCode:
		return []int{50}

	// Zone 51
	case CountryNetherlandsNewGuinea_deleted.EntityName, CountryNetherlandsNewGuinea_deleted.EntityCode,
		CountryPapuaNewGuinea.EntityName, CountryPapuaNewGuinea.EntityCode,
		CountryPapuaTerritory_deleted.EntityName, CountryPapuaTerritory_deleted.EntityCode,
		CountrySolomonIslands.EntityName, CountrySolomonIslands.EntityCode,
		CountryTemotuProvince.EntityName, CountryTemotuProvince.EntityCode,
		CountryTerritoryOfNewGuinea_deleted.EntityName, CountryTerritoryOfNewGuinea_deleted.EntityCode:
		return []int{51}

	// Zone 52
	case CountryAngola.EntityName, CountryAngola.EntityCode,
		CountryAnnobonIsland.EntityName, CountryAnnobonIsland.EntityCode,
		CountryBurundi.EntityName, CountryBurundi.EntityCode,
		CountryDemocraticRepublicOfTheCongo.EntityName, CountryDemocraticRepublicOfTheCongo.EntityCode,
		CountryGabon.EntityName, CountryGabon.EntityCode,
		CountryRepublicOfTheCongo.EntityName, CountryRepublicOfTheCongo.EntityCode,
		CountryRuandaUrundi_deleted.EntityName, CountryRuandaUrundi_deleted.EntityCode,
		CountryRwanda.EntityName, CountryRwanda.EntityCode:
		return []int{52}

	// Zone 53
	case CountryAgalegaStBrandonIslands.EntityName, CountryAgalegaStBrandonIslands.EntityCode,
		CountryAldabra_deleted.EntityName, CountryAldabra_deleted.EntityCode,
		CountryComoros.EntityName, CountryComoros.EntityCode,
		CountryComoros_deleted.EntityName, CountryComoros_deleted.EntityCode,
		CountryDesroches_deleted.EntityName, CountryDesroches_deleted.EntityCode,
		CountryFarquhar_deleted.EntityName, CountryFarquhar_deleted.EntityCode,
		CountryGeyserReef_deleted.EntityName, CountryGeyserReef_deleted.EntityCode,
		CountryGloriosoIslands.EntityName, CountryGloriosoIslands.EntityCode,
		CountryJuanDeNovaEuropa.EntityName, CountryJuanDeNovaEuropa.EntityCode,
		CountryMadagascar.EntityName, CountryMadagascar.EntityCode,
		CountryMalawi.EntityName, CountryMalawi.EntityCode,
		CountryMauritius.EntityName, CountryMauritius.EntityCode,
		CountryMayotte.EntityName, CountryMayotte.EntityCode,
		CountryMozambique.EntityName, CountryMozambique.EntityCode,
		CountryReunionIsland.EntityName, CountryReunionIsland.EntityCode,
		CountryRodriguesIsland.EntityName, CountryRodriguesIsland.EntityCode,
		CountrySeychelles.EntityName, CountrySeychelles.EntityCode,
		CountryTanzania.EntityName, CountryTanzania.EntityCode,
		CountryTromelinIsland.EntityName, CountryTromelinIsland.EntityCode,
		CountryZambia.EntityName, CountryZambia.EntityCode,
		CountryZanzibar_deleted.EntityName, CountryZanzibar_deleted.EntityCode,
		CountryZimbabwe.EntityName, CountryZimbabwe.EntityCode:
		return []int{53}

	// Zone 54
	case CountryBritishNorthBorneo_deleted.EntityName, CountryBritishNorthBorneo_deleted.EntityCode,
		CountryBruneiDarussalam.EntityName, CountryBruneiDarussalam.EntityCode,
		CountryCelebeMoluccaIslands_deleted.EntityName, CountryCelebeMoluccaIslands_deleted.EntityCode,
		CountryChristmasIsland.EntityName, CountryChristmasIsland.EntityCode,
		CountryCocosKeelingIslands.EntityName, CountryCocosKeelingIslands.EntityCode,
		CountryEastMalaysia.EntityName, CountryEastMalaysia.EntityCode,
		CountryJava_deleted.EntityName, CountryJava_deleted.EntityCode,
		CountryMalaya_deleted.EntityName, CountryMalaya_deleted.EntityCode,
		CountryNetherlandsBorneo_deleted.EntityName, CountryNetherlandsBorneo_deleted.EntityCode,
		CountryPortugueseTimor_deleted.EntityName, CountryPortugueseTimor_deleted.EntityCode,
		CountrySarawak_deleted.EntityName, CountrySarawak_deleted.EntityCode,
		CountrySingapore.EntityName, CountrySingapore.EntityCode,
		CountrySumatra_deleted.EntityName, CountrySumatra_deleted.EntityCode,
		CountryTimorLeste.EntityName, CountryTimorLeste.EntityCode,
		CountryWestMalaysia.EntityName, CountryWestMalaysia.EntityCode:
		return []int{54}

	// Zone 55
	case CountryWillisIsland.EntityName, CountryWillisIsland.EntityCode:
		return []int{55}

	// Zone 56
	case CountryChesterfieldIslands.EntityName, CountryChesterfieldIslands.EntityCode,
		CountryConwayReef.EntityName, CountryConwayReef.EntityCode,
		CountryFiji.EntityName, CountryFiji.EntityCode,
		CountryMellishReef.EntityName, CountryMellishReef.EntityCode,
		CountryNewCaledonia.EntityName, CountryNewCaledonia.EntityCode,
		CountryRotumaIsland.EntityName, CountryRotumaIsland.EntityCode,
		CountryVanuatu.EntityName, CountryVanuatu.EntityCode:
		return []int{56}

	// Zone 57
	case CountryBotswana.EntityName, CountryBotswana.EntityCode,
		CountryKingdomOfEswatini.EntityName, CountryKingdomOfEswatini.EntityCode,
		CountryLesotho.EntityName, CountryLesotho.EntityCode,
		CountryNamibia.EntityName, CountryNamibia.EntityCode,
		CountryPenguinIslands_deleted.EntityName, CountryPenguinIslands_deleted.EntityCode,
		CountryPrinceEdwardMarionIslands.EntityName, CountryPrinceEdwardMarionIslands.EntityCode,
		CountryRepublicOfSouthAfrica.EntityName, CountryRepublicOfSouthAfrica.EntityCode,
		CountryWalvisBay_deleted.EntityName, CountryWalvisBay_deleted.EntityCode:
		return []int{57}

	// Zones 58 and 59 are just Australia
	// Zone 60
	case CountryChathamIslands.EntityName, CountryChathamIslands.EntityCode,
		CountryKermadecIslands.EntityName, CountryKermadecIslands.EntityCode,
		CountryLordHoweIsland.EntityName, CountryLordHoweIsland.EntityCode,
		CountryMacquarieIsland.EntityName, CountryMacquarieIsland.EntityCode,
		CountryNewZealand.EntityName, CountryNewZealand.EntityCode,
		CountryNewZealandSubantarcticIslands.EntityName, CountryNewZealandSubantarcticIslands.EntityCode,
		CountryNorfolkIsland.EntityName, CountryNorfolkIsland.EntityCode:
		return []int{60}

	// Zone 61
	case CountryBakerHowlandIslands.EntityName, CountryBakerHowlandIslands.EntityCode,
		CountryHawaii.EntityName, CountryHawaii.EntityCode,
		CountryJohnstonIsland.EntityName, CountryJohnstonIsland.EntityCode,
		CountryKingmanReef_deleted.EntityName, CountryKingmanReef_deleted.EntityCode,
		CountryKureIsland.EntityName, CountryKureIsland.EntityCode,
		CountryMidwayIsland.EntityName, CountryMidwayIsland.EntityCode:
		return []int{61}

	// Zone 62
	case CountryAmericanSamoa.EntityName, CountryAmericanSamoa.EntityCode,
		CountryCKiribatiBritishPhoenixIslands.EntityName, CountryCKiribatiBritishPhoenixIslands.EntityCode,
		CountryMinervaReef_deleted.EntityName, CountryMinervaReef_deleted.EntityCode,
		CountryNiue.EntityName, CountryNiue.EntityCode,
		CountryNorthCookIslands.EntityName, CountryNorthCookIslands.EntityCode,
		CountrySamoa.EntityName, CountrySamoa.EntityCode,
		CountrySouthCookIslands.EntityName, CountrySouthCookIslands.EntityCode,
		CountrySwainsIsland.EntityName, CountrySwainsIsland.EntityCode,
		CountryTokelauIslands.EntityName, CountryTokelauIslands.EntityCode,
		CountryTonga.EntityName, CountryTonga.EntityCode,
		CountryWallisFutunaIslands.EntityName, CountryWallisFutunaIslands.EntityCode:
		return []int{62}

	// Zone 63
	case CountryAustralIsland.EntityName, CountryAustralIsland.EntityCode,
		CountryDucieIsland.EntityName, CountryDucieIsland.EntityCode,
		CountryEasterIsland.EntityName, CountryEasterIsland.EntityCode,
		CountryFrenchPolynesia.EntityName, CountryFrenchPolynesia.EntityCode,
		CountryMarquesasIslands.EntityName, CountryMarquesasIslands.EntityCode,
		CountryPitcairnIsland.EntityName, CountryPitcairnIsland.EntityCode:
		return []int{63}

	// Zone 64
	case CountryGuam.EntityName, CountryGuam.EntityCode,
		CountryMarianaIslands.EntityName, CountryMarianaIslands.EntityCode,
		CountryPalau.EntityName, CountryPalau.EntityCode:
		return []int{64}

	// Zone 65
	case CountryBanabaIslandOceanIsland.EntityName, CountryBanabaIslandOceanIsland.EntityCode,
		CountryMarshallIslands.EntityName, CountryMarshallIslands.EntityCode,
		CountryMicronesia.EntityName, CountryMicronesia.EntityCode,
		CountryNauru.EntityName, CountryNauru.EntityCode,
		CountryTuvalu.EntityName, CountryTuvalu.EntityCode,
		CountryWKiribatiGilbertIslands.EntityName, CountryWKiribatiGilbertIslands.EntityCode,
		CountryWakeIsland.EntityName, CountryWakeIsland.EntityCode:
		return []int{65}

	// Zone 66
	case CountryAscensionIsland.EntityName, CountryAscensionIsland.EntityCode:
		return []int{66}
	case CountryStHelena.EntityName, CountryStHelena.EntityCode:
		return []int{66}
	case CountryTristanDaCunhaGoughIsland.EntityName, CountryTristanDaCunhaGoughIsland.EntityCode:
		return []int{66}

	// Zone 67
	case CountryBouvet.EntityName, CountryBouvet.EntityCode:
		return []int{67}

	// Zone 68
	case CountryAmsterdamStPaulIslands.EntityName, CountryAmsterdamStPaulIslands.EntityCode:
		return []int{68}
	case CountryCrozetIsland.EntityName, CountryCrozetIsland.EntityCode:
		return []int{68}
	case CountryHeardIsland.EntityName, CountryHeardIsland.EntityCode:
		return []int{68}
	case CountryKerguelenIslands.EntityName, CountryKerguelenIslands.EntityCode:
		return []int{68}

	// Zones 69 through 71 are just Antarctica
	case CountryPeter1Island.EntityName, CountryPeter1Island.EntityCode:
		return []int{72}

	// Zone 73
	case CountrySouthGeorgiaIsland.EntityName, CountrySouthGeorgiaIsland.EntityCode:
		return []int{73}
	case CountrySouthOrkneyIslands.EntityName, CountrySouthOrkneyIslands.EntityCode:
		return []int{73}
	case CountrySouthSandwichIslands.EntityName, CountrySouthSandwichIslands.EntityCode:
		return []int{73}

	// Zone 74 is the South Pole
	case CountrySouthShetlandIslands.EntityName, CountrySouthShetlandIslands.EntityCode:
		return []int{73}

	// Zone 75
	case CountryFranzJosefLand.EntityName, CountryFranzJosefLand.EntityCode:
		return []int{75}

	// Zones 76 through 89 are entirely ocean
	// Zone 90
	case CountryMinamiTorishima.EntityName, CountryMinamiTorishima.EntityCode:
		return []int{90}
	}
}
