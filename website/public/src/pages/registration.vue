<template>
	<div>
		<navbar scrollable="true"></navbar>
		<div class="section-reg background--blue-10">
			<div class="container">
				<div class="reg-card center-xs">
					<div class="progress" v-if="step !== 0">
						<div class="icon icon-expand icon-rotate-90" v-show="step >= 2 && step < 4" v-on:click="stepBack"></div>
						<div v-bind:class="{'step': true, 'done': step >= 1}"></div>
						<div v-bind:class="{'step': true, 'done': step >= 2}"></div>
						<div v-bind:class="{'step': true, 'done': step >= 3}"></div>
						<div v-bind:class="{'step': true, 'done': step >= 4}"></div>
					</div>
					<registrationemail v-if="step === 0" :initialemail="email"></registrationemail>
					<registrationpassword v-if="step === 1" :code="code" v-on:success="step++"></registrationpassword>
					<registrationpersonal v-if="step === 2" :code="code" v-on:success="step++"></registrationpersonal>
					<registrationcompletion v-if="step === 3" :code="code" v-on:success="step++"></registrationcompletion>
					<registrationdone v-if="step > 3"></registrationdone>
					<registrationcancel v-if="step === -1" :code="code" v-on:cancel="loadStep"></registrationcancel>
					<div v-if="notfound">
						<h3 class="reg-headline">
							{{$t("not_found_title")}}
						</h3>
						<p>{{$t("not_found")}}</p>
						<button class="col-sm-12 reg-button--primary" v-on:click="toRegistration">{{$t("button_registration")}}</button>
					</div>
					<div class="reg-footer">
						<router-link to="/terms" target="_blank" ref="noreferrer">{{$t("link_terms")}}</router-link>
						<span class="dot">·</span>
						<router-link to="/privacy" target="_blank" ref="noreferrer">{{$t("link_privacy")}}</router-link>
						<template v-if="code && !notfound && step !== -1 && step < 4">
							<span class="dot">·</span>
							<a href="#" v-on:click.prevent="step = -1">{{$t("link_cancel_registration")}}</a>
						</template>
					</div>
				</div>
				<div class="container center-lg center-md center-sm center-xs">
					<div class="spacer-16"></div>
					<q>{{$t("account")}} <span class="color--blue-100 cursor-pointer hover-underline" v-on:click="login">{{$t("login")}}</span></q>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
	import {AuthService} from "../service";
	import {title} from "../mixins";
	import {navbar, registrationemail, registrationcancel, registrationpassword, registrationpersonal, registrationcompletion, registrationdone} from "../components";

	export default {
		mixins: [title],
		components: {
			navbar,
			registrationemail,
			registrationcancel,
			registrationpassword,
			registrationpersonal,
			registrationcompletion,
			registrationdone
		},
		data() {
			return {
				email: "",
				code: "",
				step: -2,
				notfound: false
			};
		},
		mounted() {
			this.email = this.$route.query.email;
			let code = this.$route.query.code;

			if(code) {
				this.code = code;
				this.loadStep();
			}
			else {
				this.step = 0;
			}
		},
		methods: {
			loadStep() {
				AuthService.getRegistrationStep(this.code)
				.then(step => {
					this.step = step;
				})
				.catch(e => {
					this.notfound = true;
					this.setError(e);
				});
			},
			toRegistration() {
				this.code = "";
				this.step = 0;
				this.notfound = false;
				this.$router.push("/registration");
			},
			stepBack() {
				this.step--;
			},
			login() {
				this.$store.dispatch("login");
			}
		}
	}
</script>

<style scoped>
	.reg-card {
		text-align: left;
	}
</style>

<i18n>
	{
		"en": {
			"not_found_title": "Not found",
			"not_found": "The registration could not be found. The registration email can be send again if you register again with the same email address.",
			"button_registration": "Register now",
			"link_cancel_registration": "Cancel",
			"link_terms": "Terms",
			"link_privacy": "Privacy",
			"account": "Already have an account?",
			"login": "Login"
		},
		"de": {
			"not_found_title": "Nicht gefunden ",
			"not_found": "Die Registierung konnte nicht gefunden werden. Die Registrierungs-E-Mail kann ein weiteres Mal gesendet werden, wenn du dich mit der selben E-Mail-Adresse erneut registrierst.",
			"button_registration": "Jetzt registrieren",
			"link_cancel_registration": "Abbrechen",
			"link_terms": "Nutzungsbedingungen",
			"link_privacy": "Datenschutzerklärung",
			"account": "Konto vorhanden?",
			"login": "Anmelden"
		}
	}
</i18n>
