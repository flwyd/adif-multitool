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

// CQZoneFor returns a list of CQ Zone numbers for a DXCC entity code or country
// name.  Returns an empty slice if the entity code or country is not known.
// For some countries which span multiple CQ Zones, the CqZone property of the
// SecondaryAdministrativeSubdivision enumeration indicates the station's zone.
// This list was compiled from https://cqww.com/cq_waz_list.htm which was last
// updated April 1st, 2018 from
// https://www.cq-amateur-radio.com/cq_awards/cq_waz_awards/cq_waz_list.html
// which is archived at
// https://web.archive.org/web/20240119065801/https://www.cq-amateur-radio.com/cq_awards/cq_waz_awards/cq_waz_list.html
// See also https://www.mapability.com/ei8ic/maps/cqzone.php and
// https://zone-info.eu/
// Desecheo Island and Republic of Kosovo have been added to the CQWW list.
// Deleted DXCC entities have been added based on current CQ zone boundaries.
//
// TODO Take ISO codes as well?
func CQZoneFor(s string) []int {
	switch strings.ToUpper(s) {
	default:
		return []int{}
	// Multi-zone countries
	case CountryCanada.EntityName, CountryCanada.EntityCode:
		return []int{1, 2, 3, 4, 5}
	case CountryUnitedStatesOfAmerica.EntityName, CountryUnitedStatesOfAmerica.EntityCode:
		return []int{3, 4, 5}
	case CountryAsiaticRussia.EntityName, CountryAsiaticRussia.EntityCode:
		return []int{17, 18, 19, 23}
	case CountryYemen.EntityName, CountryYemen.EntityCode:
		return []int{21, 37} // Socotra and Abd al Kuri Islands are in 37
	case CountryChina.EntityName, CountryChina.EntityCode:
		return []int{23, 24}
	case CountryAustralia.EntityName, CountryAustralia.EntityCode:
		return []int{29, 30}
	case CountryAntarctica.EntityName, CountryAntarctica.EntityCode:
		return []int{12, 13, 29, 30, 32, 28, 39}
	case CountryNewfoundlandLabrador_deleted.EntityName, CountryNewfoundlandLabrador_deleted.EntityCode:
		return []int{2, 5}

	// Within a zone, order follows CQWW list, which is by callsign prefix

	// Zone 1: Northwestern North America (also Canada)
	case CountryAlaska.EntityName, CountryAlaska.EntityCode:
		return []int{1}

	// Zone 2: Northeastern North America (just Canada)
	// Zone 3: Western North America (just USA and Canada)
	// Zone 4: Central North America (just USA and Canada)

	// Zone 5: Eastern North America (also USA and Canada)
	case CountryStPaulIsland.EntityName, CountryStPaulIsland.EntityCode,
		CountrySableIsland.EntityName, CountrySableIsland.EntityCode,
		CountryStPierreMiquelon.EntityName, CountryStPierreMiquelon.EntityCode,
		CountryBermuda.EntityName, CountryBermuda.EntityCode,
		CountryUnitedNationsHq.EntityName, CountryUnitedNationsHq.EntityCode:
		return []int{5}

	// Zone 6: Southern North America
	case CountryMexico.EntityName, CountryMexico.EntityCode,
		CountryRevillagigedo.EntityName, CountryRevillagigedo.EntityCode:
		return []int{6}

	// Zone 7: Central America
	case CountryClippertonIsland.EntityName, CountryClippertonIsland.EntityCode,
		CountrySanAndresProvidencia.EntityName, CountrySanAndresProvidencia.EntityCode,
		CountryPanama.EntityName, CountryPanama.EntityCode,
		CountryHonduras.EntityName, CountryHonduras.EntityCode,
		CountryGuatemala.EntityName, CountryGuatemala.EntityCode,
		CountryCostaRica.EntityName, CountryCostaRica.EntityCode,
		CountryCocosIsland.EntityName, CountryCocosIsland.EntityCode,
		CountryBelize.EntityName, CountryBelize.EntityCode,
		CountryNicaragua.EntityName, CountryNicaragua.EntityCode,
		CountryElSalvador.EntityName, CountryElSalvador.EntityCode,
		CountryCanalZone_deleted.EntityName, CountryCanalZone_deleted.EntityCode,
		CountrySerranaBankRoncadorCay_deleted.EntityName, CountrySerranaBankRoncadorCay_deleted.EntityCode,
		CountrySwanIslands_deleted.EntityName, CountrySwanIslands_deleted.EntityCode:
		return []int{7}

	// Zone 8: West Indies
	case CountryBahamas.EntityName, CountryBahamas.EntityCode,
		CountryCuba.EntityName, CountryCuba.EntityCode,
		CountryGuadeloupe.EntityName, CountryGuadeloupe.EntityCode,
		CountrySaintBarthelemy.EntityName, CountrySaintBarthelemy.EntityCode,
		CountryMartinique.EntityName, CountryMartinique.EntityCode,
		CountrySintMaarten.EntityName, CountrySintMaarten.EntityCode,
		CountryHaiti.EntityName, CountryHaiti.EntityCode,
		CountryDominicanRepublic.EntityName, CountryDominicanRepublic.EntityCode,
		CountryGrenada.EntityName, CountryGrenada.EntityCode,
		CountryStLucia.EntityName, CountryStLucia.EntityCode,
		CountryDominica.EntityName, CountryDominica.EntityCode,
		CountryStVincent.EntityName, CountryStVincent.EntityCode,
		CountryGuantanamoBay.EntityName, CountryGuantanamoBay.EntityCode,
		CountryNavassaIsland.EntityName, CountryNavassaIsland.EntityCode,
		CountryVirginIslands.EntityName, CountryVirginIslands.EntityCode,
		CountryPuertoRico.EntityName, CountryPuertoRico.EntityCode,
		CountryDesecheoIsland.EntityName, CountryDesecheoIsland.EntityCode,
		CountrySabaStEustatius.EntityName, CountrySabaStEustatius.EntityCode,
		CountrySaintMartin.EntityName, CountrySaintMartin.EntityCode,
		CountryAntiguaBarbuda.EntityName, CountryAntiguaBarbuda.EntityCode,
		CountryStKittsNevis.EntityName, CountryStKittsNevis.EntityCode,
		CountryAnguilla.EntityName, CountryAnguilla.EntityCode,
		CountryMontserrat.EntityName, CountryMontserrat.EntityCode,
		CountryBritishVirginIslands.EntityName, CountryBritishVirginIslands.EntityCode,
		CountryTurksCaicosIslands.EntityName, CountryTurksCaicosIslands.EntityCode,
		CountryAvesIsland.EntityName, CountryAvesIsland.EntityCode,
		CountryCaymanIslands.EntityName, CountryCaymanIslands.EntityCode,
		CountryJamaica.EntityName, CountryJamaica.EntityCode,
		CountryBarbados.EntityName, CountryBarbados.EntityCode,
		CountryBajoNuevo_deleted.EntityName, CountryBajoNuevo_deleted.EntityCode,
		CountryStMaartenSabaStEustatius_deleted.EntityName, CountryStMaartenSabaStEustatius_deleted.EntityCode:
		return []int{8}

	// Zone 9: Northern South America
	case CountryFrenchGuiana.EntityName, CountryFrenchGuiana.EntityCode,
		CountryColombia.EntityName, CountryColombia.EntityCode,
		CountryMalpeloIsland.EntityName, CountryMalpeloIsland.EntityCode,
		CountryCuracao.EntityName, CountryCuracao.EntityCode,
		CountryBonaire.EntityName, CountryBonaire.EntityCode,
		CountrySuriname.EntityName, CountrySuriname.EntityCode,
		CountryVenezuela.EntityName, CountryVenezuela.EntityCode,
		CountryGuyana.EntityName, CountryGuyana.EntityCode,
		CountryAruba.EntityName, CountryAruba.EntityCode,
		CountryTrinidadTobago.EntityName, CountryTrinidadTobago.EntityCode,
		CountryBonaireCuracao_deleted.EntityName, CountryBonaireCuracao_deleted.EntityCode:
		return []int{9}

	// Zone 10: Western South America
	case CountryBolivia.EntityName, CountryBolivia.EntityCode,
		CountryEcuador.EntityName, CountryEcuador.EntityCode,
		CountryGalapagosIslands.EntityName, CountryGalapagosIslands.EntityCode,
		CountryPeru.EntityName, CountryPeru.EntityCode:
		return []int{10}

	// Zone 11: Central South America
	case CountryBrazil.EntityName, CountryBrazil.EntityCode,
		CountryFernandoDeNoronha.EntityName, CountryFernandoDeNoronha.EntityCode,
		CountryStPeterStPaulRocks.EntityName, CountryStPeterStPaulRocks.EntityCode,
		CountryTrindadeMartimVazIslands.EntityName, CountryTrindadeMartimVazIslands.EntityCode,
		CountryParaguay.EntityName, CountryParaguay.EntityCode:
		return []int{11}

	// Zone 12: Southwest South America (also Antarctica)
	case CountryChile.EntityName, CountryChile.EntityCode,
		CountryEasterIsland.EntityName, CountryEasterIsland.EntityCode,
		CountryJuanFernandezIslands.EntityName, CountryJuanFernandezIslands.EntityCode,
		CountrySanFelixSanAmbrosio.EntityName, CountrySanFelixSanAmbrosio.EntityCode,
		CountryPeter1Island.EntityName, CountryPeter1Island.EntityCode:
		return []int{12}

	// Zone 13: Southeast South America (also Antarctica)
	case CountryUruguay.EntityName, CountryUruguay.EntityCode,
		CountryArgentina.EntityName, CountryArgentina.EntityCode,
		CountryFalklandIslands.EntityName, CountryFalklandIslands.EntityCode,
		CountrySouthGeorgiaIsland.EntityName, CountrySouthGeorgiaIsland.EntityCode,
		CountrySouthOrkneyIslands.EntityName, CountrySouthOrkneyIslands.EntityCode,
		CountrySouthShetlandIslands.EntityName, CountrySouthShetlandIslands.EntityCode,
		CountrySouthSandwichIslands.EntityName, CountrySouthSandwichIslands.EntityCode:
		return []int{13}

	// Zone 14: Western Europe
	case CountryAndorra.EntityName, CountryAndorra.EntityCode,
		CountryPortugal.EntityName, CountryPortugal.EntityCode,
		CountryAzores.EntityName, CountryAzores.EntityCode,
		CountryFederalRepublicOfGermany.EntityName, CountryFederalRepublicOfGermany.EntityCode,
		CountrySpain.EntityName, CountrySpain.EntityCode,
		CountryBalearicIslands.EntityName, CountryBalearicIslands.EntityCode,
		CountryIreland.EntityName, CountryIreland.EntityCode,
		CountryFrance.EntityName, CountryFrance.EntityCode,
		CountryEngland.EntityName, CountryEngland.EntityCode,
		CountryIsleOfMan.EntityName, CountryIsleOfMan.EntityCode,
		CountryNorthernIreland.EntityName, CountryNorthernIreland.EntityCode,
		CountryJersey.EntityName, CountryJersey.EntityCode,
		CountryScotland.EntityName, CountryScotland.EntityCode,
		CountryGuernsey.EntityName, CountryGuernsey.EntityCode,
		CountryWales.EntityName, CountryWales.EntityCode,
		CountrySwitzerland.EntityName, CountrySwitzerland.EntityCode,
		CountryLiechtenstein.EntityName, CountryLiechtenstein.EntityCode,
		CountryNorway.EntityName, CountryNorway.EntityCode,
		CountryLuxembourg.EntityName, CountryLuxembourg.EntityCode,
		CountryBelgium.EntityName, CountryBelgium.EntityCode,
		CountryFaroeIslands.EntityName, CountryFaroeIslands.EntityCode,
		CountryDenmark.EntityName, CountryDenmark.EntityCode,
		CountryNetherlands.EntityName, CountryNetherlands.EntityCode,
		CountrySweden.EntityName, CountrySweden.EntityCode,
		CountryGibraltar.EntityName, CountryGibraltar.EntityCode,
		CountryMonaco.EntityName, CountryMonaco.EntityCode,
		CountryItuHq.EntityName, CountryItuHq.EntityCode,
		CountryGermany_deleted.EntityName, CountryGermany_deleted.EntityCode,
		CountrySaar_deleted.EntityName, CountrySaar_deleted.EntityCode,
		CountryGermanDemocraticRepublic_deleted.EntityName, CountryGermanDemocraticRepublic_deleted.EntityCode:
		return []int{14}

	// Zone 15: Central Europe
	case CountryBosniaHerzegovina.EntityName, CountryBosniaHerzegovina.EntityCode,
		CountryEstonia.EntityName, CountryEstonia.EntityCode,
		CountryHungary.EntityName, CountryHungary.EntityCode,
		CountryVatican.EntityName, CountryVatican.EntityCode,
		CountryItaly.EntityName, CountryItaly.EntityCode,
		CountrySardinia.EntityName, CountrySardinia.EntityCode,
		CountryLithuania.EntityName, CountryLithuania.EntityCode,
		CountryAustria.EntityName, CountryAustria.EntityCode,
		CountryFinland.EntityName, CountryFinland.EntityCode,
		CountryAlandIslands.EntityName, CountryAlandIslands.EntityCode,
		CountryMarketReef.EntityName, CountryMarketReef.EntityCode,
		CountryCzechRepublic.EntityName, CountryCzechRepublic.EntityCode,
		CountrySlovakRepublic.EntityName, CountrySlovakRepublic.EntityCode,
		CountrySlovenia.EntityName, CountrySlovenia.EntityCode,
		CountryPoland.EntityName, CountryPoland.EntityCode,
		CountrySanMarino.EntityName, CountrySanMarino.EntityCode,
		CountryCorsica.EntityName, CountryCorsica.EntityCode,
		CountryKaliningrad.EntityName, CountryKaliningrad.EntityCode,
		CountryLatvia.EntityName, CountryLatvia.EntityCode,
		CountrySerbia.EntityName, CountrySerbia.EntityCode,
		CountryRepublicOfKosovo.EntityName, CountryRepublicOfKosovo.EntityCode,
		CountryAlbania.EntityName, CountryAlbania.EntityCode,
		CountryNorthMacedoniaRepublicOf.EntityName, CountryNorthMacedoniaRepublicOf.EntityCode,
		CountrySovereignMilitaryOrderOfMalta.EntityName, CountrySovereignMilitaryOrderOfMalta.EntityCode,
		CountryMontenegro.EntityName, CountryMontenegro.EntityCode,
		CountryCroatia.EntityName, CountryCroatia.EntityCode,
		CountryMalta.EntityName, CountryMalta.EntityCode,
		CountryCzechoslovakia_deleted.EntityName, CountryCzechoslovakia_deleted.EntityCode,
		CountryTrieste_deleted.EntityName, CountryTrieste_deleted.EntityCode:
		return []int{15}

	// Zone 16: Eastern Europe
	case CountryBelarus.EntityName, CountryBelarus.EntityCode,
		CountryMoldova.EntityName, CountryMoldova.EntityCode,
		CountryEuropeanRussia.EntityName, CountryEuropeanRussia.EntityCode,
		CountryUkraine.EntityName, CountryUkraine.EntityCode,
		CountryKareloFinnishRepublic_deleted.EntityName, CountryKareloFinnishRepublic_deleted.EntityCode,
		CountryMalyjVysotskijIsland_deleted.EntityName, CountryMalyjVysotskijIsland_deleted.EntityCode:
		return []int{16}

	// Zone 17: Western Siberia (also Asiatic Russia)
	case CountryKyrgyzstan.EntityName, CountryKyrgyzstan.EntityCode,
		CountryTajikistan.EntityName, CountryTajikistan.EntityCode,
		CountryTurkmenistan.EntityName, CountryTurkmenistan.EntityCode,
		CountryUzbekistan.EntityName, CountryUzbekistan.EntityCode,
		CountryKazakhstan.EntityName, CountryKazakhstan.EntityCode:
		return []int{17}

	// Zone 18: Central Siberia (just Asiatic Russia)
	// Zone 19: Eastern Siberia (just Asiatic Russia)

	// Zone 20: Balkan
	case CountryPalestine.EntityName, CountryPalestine.EntityCode,
		CountryJordan.EntityName, CountryJordan.EntityCode,
		CountryBulgaria.EntityName, CountryBulgaria.EntityCode,
		CountryLebanon.EntityName, CountryLebanon.EntityCode,
		CountryGreece.EntityName, CountryGreece.EntityCode,
		CountryMountAthos.EntityName, CountryMountAthos.EntityCode,
		CountryDodecanese.EntityName, CountryDodecanese.EntityCode,
		CountryCrete.EntityName, CountryCrete.EntityCode,
		CountryTurkey.EntityName, CountryTurkey.EntityCode,
		CountrySyria.EntityName, CountrySyria.EntityCode,
		CountryRomania.EntityName, CountryRomania.EntityCode,
		CountryUkSovereignBaseAreasOnCyprus.EntityName, CountryUkSovereignBaseAreasOnCyprus.EntityCode,
		CountryIsrael.EntityName, CountryIsrael.EntityCode,
		CountryCyprus.EntityName, CountryCyprus.EntityCode:
		return []int{20}

	// Zone 21: Southwestern Asia (also Yemen)
	case CountryOman.EntityName, CountryOman.EntityCode,
		CountryUnitedArabEmirates.EntityName, CountryUnitedArabEmirates.EntityCode,
		CountryQatar.EntityName, CountryQatar.EntityCode,
		CountryBahrain.EntityName, CountryBahrain.EntityCode,
		CountryPakistan.EntityName, CountryPakistan.EntityCode,
		CountryArmenia.EntityName, CountryArmenia.EntityCode,
		CountryIran.EntityName, CountryIran.EntityCode,
		CountrySaudiArabia.EntityName, CountrySaudiArabia.EntityCode,
		CountryAfghanistan.EntityName, CountryAfghanistan.EntityCode,
		CountryIraq.EntityName, CountryIraq.EntityCode,
		CountryAzerbaijan.EntityName, CountryAzerbaijan.EntityCode,
		CountryGeorgia.EntityName, CountryGeorgia.EntityCode,
		CountryKuwait.EntityName, CountryKuwait.EntityCode,
		CountryAbuAilIslands_deleted.EntityName, CountryAbuAilIslands_deleted.EntityCode,
		CountryKuwaitSaudiArabiaNeutralZone_deleted.EntityName, CountryKuwaitSaudiArabiaNeutralZone_deleted.EntityCode,
		CountryKamaranIslands_deleted.EntityName, CountryKamaranIslands_deleted.EntityCode,
		CountryYemenArabRepublic_deleted.EntityName, CountryYemenArabRepublic_deleted.EntityCode,
		CountryKuriaMuriaIsland_deleted.EntityName, CountryKuriaMuriaIsland_deleted.EntityCode,
		CountrySaudiArabiaIraqNeutralZone_deleted.EntityName, CountrySaudiArabiaIraqNeutralZone_deleted.EntityCode,
		CountryPeoplesDemocraticRepOfYemen_deleted.EntityName, CountryPeoplesDemocraticRepOfYemen_deleted.EntityCode:
		return []int{21}

	// Zone 22: Southern Asia
	case CountryBhutan.EntityName, CountryBhutan.EntityCode,
		CountryBangladesh.EntityName, CountryBangladesh.EntityCode,
		CountryIndia.EntityName, CountryIndia.EntityCode,
		CountryLakshadweepIslands.EntityName, CountryLakshadweepIslands.EntityCode,
		CountrySriLanka.EntityName, CountrySriLanka.EntityCode,
		CountryMaldives.EntityName, CountryMaldives.EntityCode,
		CountryNepal.EntityName, CountryNepal.EntityCode,
		CountryDamaoDiu_deleted.EntityName, CountryDamaoDiu_deleted.EntityCode,
		CountryFrenchIndia_deleted.EntityName, CountryFrenchIndia_deleted.EntityCode,
		CountryGoa_deleted.EntityName, CountryGoa_deleted.EntityCode,
		CountrySikkim_deleted.EntityName, CountrySikkim_deleted.EntityCode:
		return []int{22}

	// Zone 23: Central Asia (also China and Asiatic Russia)
	case CountryMongolia.EntityName, CountryMongolia.EntityCode,
		CountryTibet_deleted.EntityName, CountryTibet_deleted.EntityCode:
		return []int{23}

	// Zone 24: Eastern Asia (also China)
	case CountryPratasIsland.EntityName, CountryPratasIsland.EntityCode,
		CountryTaiwan.EntityName, CountryTaiwan.EntityCode,
		CountryHongKong.EntityName, CountryHongKong.EntityCode,
		CountryMacao.EntityName, CountryMacao.EntityCode,
		CountryManchuria_deleted.EntityName, CountryManchuria_deleted.EntityCode:
		return []int{24}

	// Zone 25: Japanese
	case CountryRepublicOfKorea.EntityName, CountryRepublicOfKorea.EntityCode,
		CountryJapan.EntityName, CountryJapan.EntityCode,
		CountryDemocraticPeoplesRepOfKorea.EntityName, CountryDemocraticPeoplesRepOfKorea.EntityCode,
		CountryOkinawaRyukyuIslands_deleted.EntityName, CountryOkinawaRyukyuIslands_deleted.EntityCode:
		return []int{25}

	// Zone 26: Southeastern Asia
	case CountryVietNam.EntityName, CountryVietNam.EntityCode,
		CountryThailand.EntityName, CountryThailand.EntityCode,
		CountryAndamanNicobarIslands.EntityName, CountryAndamanNicobarIslands.EntityCode,
		CountryCambodia.EntityName, CountryCambodia.EntityCode,
		CountryLaos.EntityName, CountryLaos.EntityCode,
		CountryMyanmar.EntityName, CountryMyanmar.EntityCode,
		CountrySpratlyIslands.EntityName, CountrySpratlyIslands.EntityCode,
		CountryFrenchIndoChina_deleted.EntityName, CountryFrenchIndoChina_deleted.EntityCode:
		return []int{26}

	// Zone 27: Philippine
	case CountryScarboroughReef.EntityName, CountryScarboroughReef.EntityCode,
		CountryPhilippines.EntityName, CountryPhilippines.EntityCode,
		CountryMinamiTorishima.EntityName, CountryMinamiTorishima.EntityCode,
		CountryOgasawara.EntityName, CountryOgasawara.EntityCode,
		CountryPalau.EntityName, CountryPalau.EntityCode,
		CountryGuam.EntityName, CountryGuam.EntityCode,
		CountryMarianaIslands.EntityName, CountryMarianaIslands.EntityCode,
		CountryMicronesia.EntityName, CountryMicronesia.EntityCode,
		CountryOkinoToriShima_deleted.EntityName, CountryOkinoToriShima_deleted.EntityCode:
		return []int{27}

	// Zone 28: Indonesian
	case CountrySolomonIslands.EntityName, CountrySolomonIslands.EntityCode,
		CountryPapuaNewGuinea.EntityName, CountryPapuaNewGuinea.EntityCode,
		CountryBruneiDarussalam.EntityName, CountryBruneiDarussalam.EntityCode,
		CountryIndonesia.EntityName, CountryIndonesia.EntityCode,
		CountryTimorLeste.EntityName, CountryTimorLeste.EntityCode,
		CountryWestMalaysia.EntityName, CountryWestMalaysia.EntityCode,
		CountryEastMalaysia.EntityName, CountryEastMalaysia.EntityCode,
		CountrySingapore.EntityName, CountrySingapore.EntityCode,
		CountryBritishNorthBorneo_deleted.EntityName, CountryBritishNorthBorneo_deleted.EntityCode,
		CountryCelebeMoluccaIslands_deleted.EntityName, CountryCelebeMoluccaIslands_deleted.EntityCode,
		CountryJava_deleted.EntityName, CountryJava_deleted.EntityCode,
		CountryMalaya_deleted.EntityName, CountryMalaya_deleted.EntityCode,
		CountryNetherlandsBorneo_deleted.EntityName, CountryNetherlandsBorneo_deleted.EntityCode,
		CountryNetherlandsNewGuinea_deleted.EntityName, CountryNetherlandsNewGuinea_deleted.EntityCode,
		CountryPapuaTerritory_deleted.EntityName, CountryPapuaTerritory_deleted.EntityCode,
		CountryPortugueseTimor_deleted.EntityName, CountryPortugueseTimor_deleted.EntityCode,
		CountrySarawak_deleted.EntityName, CountrySarawak_deleted.EntityCode,
		CountrySumatra_deleted.EntityName, CountrySumatra_deleted.EntityCode,
		CountryTerritoryOfNewGuinea_deleted.EntityName, CountryTerritoryOfNewGuinea_deleted.EntityCode:
		return []int{28}

	// Zone 29: Western Australia (also Australia and Antarctica)
	case CountryChristmasIsland.EntityName, CountryChristmasIsland.EntityCode,
		CountryCocosKeelingIslands.EntityName, CountryCocosKeelingIslands.EntityCode:
		return []int{29}

	// Zone 30: Eastern Australia (also Australia and Antarctica)
	case CountryChesterfieldIslands.EntityName, CountryChesterfieldIslands.EntityCode,
		CountryLordHoweIsland.EntityName, CountryLordHoweIsland.EntityCode,
		CountryMellishReef.EntityName, CountryMellishReef.EntityCode,
		CountryWillisIsland.EntityName, CountryWillisIsland.EntityCode,
		CountryMacquarieIsland.EntityName, CountryMacquarieIsland.EntityCode:
		return []int{30}

	// Zone 31: Central Pacific
	case CountryNauru.EntityName, CountryNauru.EntityCode,
		CountryMarquesasIslands.EntityName, CountryMarquesasIslands.EntityCode,
		CountryBakerHowlandIslands.EntityName, CountryBakerHowlandIslands.EntityCode,
		CountryJohnstonIsland.EntityName, CountryJohnstonIsland.EntityCode,
		CountryMidwayIsland.EntityName, CountryMidwayIsland.EntityCode,
		CountryPalmyraJarvisIslands.EntityName, CountryPalmyraJarvisIslands.EntityCode,
		CountryHawaii.EntityName, CountryHawaii.EntityCode,
		CountryKureIsland.EntityName, CountryKureIsland.EntityCode,
		CountryWakeIsland.EntityName, CountryWakeIsland.EntityCode,
		CountryTuvalu.EntityName, CountryTuvalu.EntityCode,
		CountryWKiribatiGilbertIslands.EntityName, CountryWKiribatiGilbertIslands.EntityCode,
		CountryCKiribatiBritishPhoenixIslands.EntityName, CountryCKiribatiBritishPhoenixIslands.EntityCode,
		CountryEKiribatiLineIslands.EntityName, CountryEKiribatiLineIslands.EntityCode,
		CountryBanabaIslandOceanIsland.EntityName, CountryBanabaIslandOceanIsland.EntityCode,
		CountryMarshallIslands.EntityName, CountryMarshallIslands.EntityCode,
		CountryTokelauIslands.EntityName, CountryTokelauIslands.EntityCode,
		CountryKingmanReef_deleted.EntityName, CountryKingmanReef_deleted.EntityCode:
		return []int{31}

	// Zone 32: New Zealand
	case CountryTonga.EntityName, CountryTonga.EntityCode,
		CountryNorthCookIslands.EntityName, CountryNorthCookIslands.EntityCode,
		CountrySouthCookIslands.EntityName, CountrySouthCookIslands.EntityCode,
		CountryNewCaledonia.EntityName, CountryNewCaledonia.EntityCode,
		CountryFrenchPolynesia.EntityName, CountryFrenchPolynesia.EntityCode,
		CountryAustralIsland.EntityName, CountryAustralIsland.EntityCode,
		CountryWallisFutunaIslands.EntityName, CountryWallisFutunaIslands.EntityCode,
		CountryTemotuProvince.EntityName, CountryTemotuProvince.EntityCode,
		CountryAmericanSamoa.EntityName, CountryAmericanSamoa.EntityCode,
		CountrySwainsIsland.EntityName, CountrySwainsIsland.EntityCode,
		CountryNorfolkIsland.EntityName, CountryNorfolkIsland.EntityCode,
		CountryPitcairnIsland.EntityName, CountryPitcairnIsland.EntityCode,
		CountryDucieIsland.EntityName, CountryDucieIsland.EntityCode,
		CountryVanuatu.EntityName, CountryVanuatu.EntityCode,
		CountryNiue.EntityName, CountryNiue.EntityCode,
		CountryNewZealand.EntityName, CountryNewZealand.EntityCode,
		CountryNewZealandSubantarcticIslands.EntityName, CountryNewZealandSubantarcticIslands.EntityCode,
		CountryChathamIslands.EntityName, CountryChathamIslands.EntityCode,
		CountryKermadecIslands.EntityName, CountryKermadecIslands.EntityCode,
		CountryFiji.EntityName, CountryFiji.EntityCode,
		CountryRotumaIsland.EntityName, CountryRotumaIsland.EntityCode,
		CountryConwayReef.EntityName, CountryConwayReef.EntityCode,
		CountrySamoa.EntityName, CountrySamoa.EntityCode,
		CountryMinervaReef_deleted.EntityName, CountryMinervaReef_deleted.EntityCode:
		return []int{32}

	// Zone 33: Northwestern Africa
	case CountryMorocco.EntityName, CountryMorocco.EntityCode,
		CountryMadeiraIslands.EntityName, CountryMadeiraIslands.EntityCode,
		CountryCanaryIslands.EntityName, CountryCanaryIslands.EntityCode,
		CountryCeutaMelilla.EntityName, CountryCeutaMelilla.EntityCode,
		CountryWesternSahara.EntityName, CountryWesternSahara.EntityCode,
		CountryTunisia.EntityName, CountryTunisia.EntityCode,
		CountryAlgeria.EntityName, CountryAlgeria.EntityCode,
		CountryIfni_deleted.EntityName, CountryIfni_deleted.EntityCode,
		CountryTangier_deleted.EntityName, CountryTangier_deleted.EntityCode:
		return []int{33}

	// Zone 34: Northeastern Africa
	case CountrySudan.EntityName, CountrySudan.EntityCode,
		CountryEgypt.EntityName, CountryEgypt.EntityCode,
		CountrySouthSudanRepublicOf.EntityName, CountrySouthSudanRepublicOf.EntityCode,
		CountryLibya.EntityName, CountryLibya.EntityCode,
		CountrySouthernSudan_deleted.EntityName, CountrySouthernSudan_deleted.EntityCode:
		return []int{34}

	// Zone 35: Central Africa
	case CountryTheGambia.EntityName, CountryTheGambia.EntityCode,
		CountryCapeVerde.EntityName, CountryCapeVerde.EntityCode,
		CountryLiberia.EntityName, CountryLiberia.EntityCode,
		CountryGuineaBissau.EntityName, CountryGuineaBissau.EntityCode,
		CountryCoteDIvoire.EntityName, CountryCoteDIvoire.EntityCode,
		CountryBenin.EntityName, CountryBenin.EntityCode,
		CountryMali.EntityName, CountryMali.EntityCode,
		CountryBurkinaFaso.EntityName, CountryBurkinaFaso.EntityCode,
		CountryGuinea.EntityName, CountryGuinea.EntityCode,
		CountryNigeria.EntityName, CountryNigeria.EntityCode,
		CountryMauritania.EntityName, CountryMauritania.EntityCode,
		CountryNiger.EntityName, CountryNiger.EntityCode,
		CountryTogo.EntityName, CountryTogo.EntityCode,
		CountrySenegal.EntityName, CountrySenegal.EntityCode,
		CountryGhana.EntityName, CountryGhana.EntityCode,
		CountrySierraLeone.EntityName, CountrySierraLeone.EntityCode,
		CountryFrenchWestAfrica_deleted.EntityName, CountryFrenchWestAfrica_deleted.EntityCode,
		CountryGoldCoastTogoland_deleted.EntityName, CountryGoldCoastTogoland_deleted.EntityCode:
		return []int{35}

	// Zone 36: Equatorial Africa
	case CountryAngola.EntityName, CountryAngola.EntityCode,
		CountrySaoTomePrincipe.EntityName, CountrySaoTomePrincipe.EntityCode,
		CountryCameroon.EntityName, CountryCameroon.EntityCode,
		CountryCentralAfrica.EntityName, CountryCentralAfrica.EntityCode,
		CountryRepublicOfTheCongo.EntityName, CountryRepublicOfTheCongo.EntityCode,
		CountryGabon.EntityName, CountryGabon.EntityCode,
		CountryChad.EntityName, CountryChad.EntityCode,
		CountryStHelena.EntityName, CountryStHelena.EntityCode,
		CountryAscensionIsland.EntityName, CountryAscensionIsland.EntityCode,
		CountryEquatorialGuinea.EntityName, CountryEquatorialGuinea.EntityCode,
		CountryAnnobonIsland.EntityName, CountryAnnobonIsland.EntityCode,
		CountryZambia.EntityName, CountryZambia.EntityCode,
		CountryDemocraticRepublicOfTheCongo.EntityName, CountryDemocraticRepublicOfTheCongo.EntityCode,
		CountryBurundi.EntityName, CountryBurundi.EntityCode,
		CountryRwanda.EntityName, CountryRwanda.EntityCode,
		CountryFrenchEquatorialAfrica_deleted.EntityName, CountryFrenchEquatorialAfrica_deleted.EntityCode,
		CountryRuandaUrundi_deleted.EntityName, CountryRuandaUrundi_deleted.EntityCode:
		return []int{36}

	// Zone 37: Eastern Africa (also Yemen)
	case CountryMozambique.EntityName, CountryMozambique.EntityCode,
		CountryEthiopia.EntityName, CountryEthiopia.EntityCode,
		CountryEritrea.EntityName, CountryEritrea.EntityCode,
		CountryDjibouti.EntityName, CountryDjibouti.EntityCode,
		CountrySomalia.EntityName, CountrySomalia.EntityCode,
		CountryTanzania.EntityName, CountryTanzania.EntityCode,
		CountryUganda.EntityName, CountryUganda.EntityCode,
		CountryKenya.EntityName, CountryKenya.EntityCode,
		CountryMalawi.EntityName, CountryMalawi.EntityCode,
		CountryBritishSomaliland_deleted.EntityName, CountryBritishSomaliland_deleted.EntityCode,
		CountryItalianSomaliland_deleted.EntityName, CountryItalianSomaliland_deleted.EntityCode,
		CountryZanzibar_deleted.EntityName, CountryZanzibar_deleted.EntityCode:
		return []int{37}

	// Zone 38: South Africa (also Antarctica)
	case CountryBotswana.EntityName, CountryBotswana.EntityCode,
		CountryNamibia.EntityName, CountryNamibia.EntityCode,
		CountryTristanDaCunhaGoughIsland.EntityName, CountryTristanDaCunhaGoughIsland.EntityCode,
		CountryZimbabwe.EntityName, CountryZimbabwe.EntityCode,
		CountryRepublicOfSouthAfrica.EntityName, CountryRepublicOfSouthAfrica.EntityCode,
		CountryPrinceEdwardMarionIslands.EntityName, CountryPrinceEdwardMarionIslands.EntityCode,
		CountryKingdomOfEswatini.EntityName, CountryKingdomOfEswatini.EntityCode,
		CountryBouvet.EntityName, CountryBouvet.EntityCode,
		CountryLesotho.EntityName, CountryLesotho.EntityCode,
		CountryWalvisBay_deleted.EntityName, CountryWalvisBay_deleted.EntityCode,
		CountryPenguinIslands_deleted.EntityName, CountryPenguinIslands_deleted.EntityCode:
		return []int{38}

	// Zone 39: Madagascar (also Antarctica)
	case CountryComoros.EntityName, CountryComoros.EntityCode,
		CountryMayotte.EntityName, CountryMayotte.EntityCode,
		CountryReunionIsland.EntityName, CountryReunionIsland.EntityCode,
		CountryGloriosoIslands.EntityName, CountryGloriosoIslands.EntityCode,
		CountryJuanDeNovaEuropa.EntityName, CountryJuanDeNovaEuropa.EntityCode,
		CountryTromelinIsland.EntityName, CountryTromelinIsland.EntityCode,
		CountryCrozetIsland.EntityName, CountryCrozetIsland.EntityCode,
		CountryKerguelenIslands.EntityName, CountryKerguelenIslands.EntityCode,
		CountryAmsterdamStPaulIslands.EntityName, CountryAmsterdamStPaulIslands.EntityCode,
		CountrySeychelles.EntityName, CountrySeychelles.EntityCode,
		CountryHeardIsland.EntityName, CountryHeardIsland.EntityCode,
		CountryChagosIslands.EntityName, CountryChagosIslands.EntityCode,
		CountryAgalegaStBrandonIslands.EntityName, CountryAgalegaStBrandonIslands.EntityCode,
		CountryMauritius.EntityName, CountryMauritius.EntityCode,
		CountryRodriguesIsland.EntityName, CountryRodriguesIsland.EntityCode,
		CountryMadagascar.EntityName, CountryMadagascar.EntityCode,
		CountryAldabra_deleted.EntityName, CountryAldabra_deleted.EntityCode,
		CountryBlenheimReef_deleted.EntityName, CountryBlenheimReef_deleted.EntityCode,
		CountryDesroches_deleted.EntityName, CountryDesroches_deleted.EntityCode,
		CountryFarquhar_deleted.EntityName, CountryFarquhar_deleted.EntityCode,
		CountryGeyserReef_deleted.EntityName, CountryGeyserReef_deleted.EntityCode:
		return []int{39}

	// Zone 40: North Atlantic
	case CountrySvalbard.EntityName, CountrySvalbard.EntityCode,
		CountryJanMayen.EntityName, CountryJanMayen.EntityCode,
		CountryGreenland.EntityName, CountryGreenland.EntityCode,
		CountryFranzJosefLand.EntityName, CountryFranzJosefLand.EntityCode,
		CountryIceland.EntityName, CountryIceland.EntityCode:
		return []int{40}
	}
}
