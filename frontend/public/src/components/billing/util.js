import {getCodeList} from "country-list";

const iso31661alpha2EU = [
    "be",
    "bg",
    "dk",
    "de",
    "ee",
    "fi",
    "fr",
    "gr",
    "ie",
    "it",
    "hr",
    "lv",
    "lt",
    "lu",
    "mt",
    "nl",
    "at",
    "pl",
    "pt",
    "ro",
    "sk",
    "si",
    "es",
    "se",
    "cz",
    "hu",
    "cy"
];

const pricePerUser = {
    yearly: 5*12,
    monthly: 6
};

const taxRateDE = 0.16;

export function getCountries() {
    let countries = getCodeList();
    let options = [];

    for(let code of Object.keys(countries)) {
        options.push({value: code.toLowerCase(), label: countries[code]});
    }

    return options;
}

export function calculateSum(plan, billableMemberCount, country, taxNumber) {
    country = country.toLowerCase();
    let sum = pricePerUser[plan]*billableMemberCount;
    sum *= 1+getTaxRate(country, taxNumber);
    return sum.toFixed(2);
}

export function getTaxRate(country, taxNumber) {
    if(country === "de" || (!taxNumber.length && isEUCountry(country))) {
        return taxRateDE;
    }

    return 0;
}

function isEUCountry(country) {
    for(let i = 0; i < iso31661alpha2EU.length; i++) {
        if(country === iso31661alpha2EU[i]) {
            return true;
        }
    }

    return false;
}
