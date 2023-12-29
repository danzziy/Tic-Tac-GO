<template>
    <button @click="onClick()" :style="{background: color}">{{ text }}</button>
</template>

<script>
export default {
    name: 'Button',
    props: {
        text: String,
        color: String,
    },
    methods: {
        onClick() {
            console.log('click')
        }
    }
}
</script>

