<template>
    <div>
        <span class="display-flex align-items-center orange-100 cursor-default" v-show="connectionState === 0">
            <i class="icon icon-sync icon-is-spinning size-40"></i>
            <p>{{$t("connecting")}}</p>
        </span>
        <span class="display-flex align-items-center green-100 cursor-default" v-show="connectionState === 1">
            <i class="icon icon-connect size-40"></i>
            <p>{{$t("connected")}}</p>
        </span>
        <span class="display-flex align-items-center red-100 cursor-default" v-show="connectionState === 2">
            <i class="icon icon-disconnect size-40"></i>
            <p>{{$t("disconnected")}}</p>
        </span>
    </div>
</template>

<script>
    export default {
        props: ["connectionState"]
    }
</script>

<i18n>
{
	"en": {
		"connected": "Connected",
		"connecting": "Connecting...",
		"disconnected": "Disconnected"
	},
	"de": {
		"connected": "Verbunden",
		"connecting": "Verbinden...",
		"disconnected": "Nicht Verbunden"
	}
}
</i18n>
