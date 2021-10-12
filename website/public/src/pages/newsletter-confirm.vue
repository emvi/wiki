<template>
    <div>
		<navbar scrollable="true"></navbar>
		<div class="section-reg background--blue-10">
			<div class="container">
				<div class="reg-card center-xs">
					<template v-if="success">
						<h3 class="h3">
							{{$t("title_success")}}
						</h3>
                        <div class="reg-icon">
		                    <div class="icon icon-mail"></div>
		                </div>
						<p>{{$t("text_success")}}</p>
                        <div class="spacer-16"></div>
                        <a href="/" class="col-sm-12 reg-button--primary">{{$t("link_home")}}</a>
					</template>
                    <template v-if="error">
						<h3 class="reg-headline">
							{{$t("title_error")}}
						</h3>
						<p>{{$t("text_error")}}</p>
                        <div class="spacer-16"></div>
                        <a href="/" class="col-sm-12 reg-button--primary">{{$t("link_home")}}</a>
					</template>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
	import {NewsletterService} from "../service";
	import {title} from "../mixins";
	import {navbar} from "../components";

	export default {
		mixins: [title],
		components: {navbar},
		data() {
			return {
				success: false,
				error: false
			};
		},
		mounted() {
			let code = this.$route.query.code;
			this.confirmNewsletter(code);
		},
		methods: {
			confirmNewsletter(code) {
				NewsletterService.confirm(code)
				.then(() => {
					this.success = true;
				})
				.catch(() => {
					this.error = true;
				});
			}
		}
	}
</script>

<i18n>
	{
		"en": {
			"title_success": "Success",
			"text_success": "Thank you for subscribing to our newsletter.",
			"title_error": "Error",
			"text_error": "We could not confirm your email address.",
			"link_home": "Back to Home"
		},
		"de": {
			"title_success": "Erfolg",
			"text_success": "Danke, dass du unseren Newsletter abonniert hast.",
			"title_error": "Fehler",
			"text_error": "Wir konnten deine E-Mail-Adresse nicht bestätigen.",
			"link_home": "Zurück zur Startseite"
		}
	}
</i18n>
