<template>
    <div class="card cursor-pointer no-select" v-on:click="open">
        <div class="card-content">
            <div class="avatar size-40">
                <i class="icon icon-billing blue-100 bg-blue-10"></i>
            </div>
            <div class="item">
                <p><strong>{{invoice.number}}</strong></p>
                <small>
                    <span>{{invoice.created | moment('LL')}}</span>
                    <span class="dot">·</span>
                    <span>{{total}}</span>
                    <span class="dot">·</span>
                    <span v-if="invoice.paid" class="green-100">{{$t("paid")}}</span>
                    <span v-else class="orange-100">{{$t("unpaid")}}</span>
                </small>
            </div>
            <emvi-avatar icon="download" size="40"></emvi-avatar>
        </div>
    </div>
</template>

<script>
    import emviAvatar from "../cmd/content/avatar.vue";

    export default {
        components: {emviAvatar},
        props: ["invoice"],
        computed: {
            total() {
                let total = this.invoice.total || 0;
                return `$${(total/100).toFixed(2)}`;
            }
        },
        methods: {
            open() {
                window.open(this.invoice.pdf_link, "_blank");
            }
        }
    }
</script>

<i18n>
    {
        "en": {
            "paid": "Paid",
            "unpaid": "Open"
        },
        "de": {
            "paid": "Bezahlt",
            "unpaid": "Offen"
        }
    }
</i18n>
