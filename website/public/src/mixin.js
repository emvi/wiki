import Vue from "vue";
import {isEmptyObject} from "./util/object";
import {ErrorService} from "./service";

Vue.mixin({
    data() {
        return {
            validation: {},
            err: null
        };
    },
    methods: {
        setError(e) {
            if(!e) {
                return;
            }

            if(!isEmptyObject(e.validation)) {
                this.validation = ErrorService.map(e.validation);
            }
            else if(e.errors && e.errors.length) {
                // display the first error only
                this.err = ErrorService.mapError(e.errors[0].message);
            }
            else if(e.exceptions && e.exceptions.length) {
                // display the first exception only
                this.err = ErrorService.mapError(e.exceptions[0]);
            }
            else {
                this.err = e;
            }
        }
    }
});
