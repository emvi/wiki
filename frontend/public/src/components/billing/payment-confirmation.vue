<template>
    <div v-if="error">
        {{error}}
    </div>
</template>

<script>
    import {mapGetters} from "vuex";
    import {BillingService} from "../../service";

    export default {
        data() {
            return {
                error: ""
            };
        },
        computed: {
            ...mapGetters(["organization"])
        },
        mounted() {
            this.confirmPaymentIfRequired();
        },
        methods: {
            confirmPaymentIfRequired() {
                if(this.organization.stripe_payment_intent_client_secret) {
                    let stripe = Stripe(EMVI_WIKI_STRIPE_PUBLIC_KEY);
                    stripe.confirmCardPayment(this.organization.stripe_payment_intent_client_secret, {
                        payment_method: this.organization.stripe_payment_method_id
                    })
                        .then((result) => {
                            if (result.error) {
                                this.error = result.error.message;
                            } else if (result.paymentIntent.status === "succeeded") {
                                BillingService.removePaymentIntentClientSecret()
                                    .then(() => {
                                        this.$store.dispatch("reload");
                                        // TODO show success
                                    })
                                    .catch(e => {
                                        console.error(e);
                                        this.setError(e);
                                    });
                            }
                        });
                }
            }
        }
    }
</script>
