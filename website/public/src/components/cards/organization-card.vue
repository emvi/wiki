<template>
	<div class="large-card">
		<div class="large-card--pic">
			<a :href="visit">
				<i class="icon icon-organization" v-if="!orga.picture"></i>
				<img v-bind:src="orga.picture" alt="" v-if="orga.picture" />
			</a>
		</div>
		<div class="large-card--content">
			<a :href="visit" class="large-card--content--title">
				{{orga.name}}
			</a>
			<div class="large-card--content--info">
				<div class="large-card--content--info--level--entry" v-if="!orga.expert">
					Entry
				</div>
				<div class="large-card--content--info--level--expert" v-if="orga.expert">
					Expert
				</div>
				<div class="large-card--content--info--with-icon">
					<i class="icon icon-user large-card--content--info--with-icon--icon"></i>
					<div class="large-card--content--info--with-icon--text">
						{{orga.member_count}}
					</div>
				</div>
				<div class="large-card--content--info--with-icon">
					<i class="icon icon-article large-card--content--info--with-icon--icon"></i>
					<div class="large-card--content--info--with-icon--text">
						{{orga.article_count}}
					</div>
				</div>
				<div class="large-card--content--info--with-icon" v-if="$mq > 512">
					<i class="icon icon-link large-card--content--info--with-icon--icon"></i>
					<a :href="visit" class="large-card--content--info--with-icon--text">
						{{orga.name_normalized}}.emvi.com
					</a>
				</div>
			</div>
		</div>
	</div>
</template>

<script>
	export default {
		props: ["orga"],
		computed: {
			visit() {
				let url = EMVI_WIKI_FRONTEND_HOST.split("/");
				return `${url[0]}//${this.orga.name_normalized}.${url[2]}`;
			}
		}
	}
</script>
