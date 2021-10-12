<template>
	<layout>
		<template>
			<notfound>{{$t("title")}}</notfound>
		</template>
	</layout>
</template>

<script>
	import {title} from "../mixins";
	import {layout, notfound} from "../components";

	export default {
		mixins: [title],
		components: {
			layout,
			notfound
		}
	}
</script>

<i18n>
	{
		"en": {
			"title": "The page does not exist."
		},
		"de": {
			"title": "Die Seite existiert nicht."
		}
	}
</i18n>
