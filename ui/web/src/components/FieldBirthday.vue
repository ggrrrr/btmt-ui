<template>
    <v-container>
        <div>
            {{ formatDoB() }}
        </div>
    </v-container>
</template>

<script setup>
import { Dob } from '@/store/people';

const props = defineProps(
    {
        dob: {
            type: Dob
        }
    }
)

function formatDoB() {
    let year = "____"
    let day = "__"
    let month = "__"
    if (props.dob.day != undefined && props.dob.day !== "") {
        day = props.dob.day
    }
    if (props.dob.month != undefined && props.dob.month !== "") {
        month = props.dob.month
    }
    if (props.dob.year != undefined && props.dob.year !== "") {
        year = props.dob.year
    }
    return `${day} / ${month} / ${year}`
}

</script>