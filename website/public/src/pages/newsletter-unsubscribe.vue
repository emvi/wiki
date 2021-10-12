<template>
    <div>
        <navbar scrollable="true"></navbar>
		<div class="section-reg background--blue-10">
			<div class="container">
				<div class="reg-card center-xs">
					<template v-if="success">
						<h3 class="reg-headline">
							{{$t("title_success")}}
						</h3>
                        <div class="reg-icon">
		                    <div class="icon icon-mail"></div>
		                </div>
						<p>{{$t("text_success")}}</p>
                        <div class="spacer-16"></div>
                        <router-link to="/" class="col-sm-12 reg-button--primary">{{$t("link_home")}}</router-link>
					</template>
                    <template v-if="error">
						<h3 class="reg-headline">
							{{$t("title_error")}}
						</h3>
						<p>{{$t("text_error")}}</p>
                        <div class="spacer-16"></div>
                        <router-link to="/" class="col-sm-12 reg-button--primary">{{$t("link_home")}}</router-link>
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
			this.unsubNewsletter(code);
		},
		methods: {
			unsubNewsletter(code) {
				NewsletterService.unsubscribe(code)
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
			"text_success": "You are no longer subscribed to our newsletter.",
			"title_error": "Error",
			"text_error": "We could not find your email address, it might be deleted already.",
			"link_home": "Back to Home"
		},
		"de": {
			"title_success": "Erfolg",
			"text_success": "Du erhälst nun keine Newsletter mehr.",
			"title_error": "Fehler",
			"text_error": "Wir konnten deine E-Mail-Adresse nicht finden, sie wurde möglicherweise bereits gelöscht.",
			"link_home": "Zurück zur Startseite"
		}
	}
</i18n>
