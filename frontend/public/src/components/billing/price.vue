<template>
    <div :class="{'plan-content': true, 'active': value === plan}">
        <input type="radio" :id="plan" name="plan"
            :value="plan"
            :checked="value === plan"
            v-on:input="$emit('input', $event.target.value)" />
        <label :for="plan">
            <p><strong>{{$t(plan)}}</strong></p>
            <span>
                <h2>{{plan === "monthly" ? "$6" : "$5"}}</h2>
                <small>{{$t("per_user_month")}}</small>
            </span>
        </label>
    </div>
</template>

<script>
    export default {
        props: ["plan", "value"]
    }
</script>

<i18n>
    {
        "en": {
            "monthly": "Bill monthly",
            "yearly": "Bill yearly",
            "per_user_month": "per user / month"
        },
        "de": {
            "monthly": "Monatliche Abrechnung",
            "yearly": "Jährliche Abrechnung",
            "per_user_month": "pro Nutzer / Monat"
        }
    }
</i18n>
