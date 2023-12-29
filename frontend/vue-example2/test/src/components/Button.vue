<template>
    <div :style="{background: color}">{{ text }}</div>
</template>

<script>
export default {
    name: 'Button',
    props: {
        text: String,
        color: String,
    },
}
</script>

<style scoped>
div {
    width: 520px;
    height: 40px;
    display: inline-flex;
    translate: 500px;

    font-family: '8bit';
    font-size: 36px;
    justify-content: center;
    color: white;
    
    cursor: pointer;

    transition: font-size 0.6s ease, font-weight 1s ease; 
}
div:hover {
    font-size: 52px; /* Larger font size when hovered */

}
</style>