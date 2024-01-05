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
    height: fit-content;
    width: fit-content;
    padding: 25px;
    margin: 10px auto;

    font-family: '8bit';
    font-size: 36px;
    justify-content: center;
    color: white;
    
    cursor: pointer;

    pointer-events: all;

    transition: font-size 0.6s ease, font-weight 1s ease; 
    top: 50%;
    left: 50%;
    transform: translate(0%, 50%); 
}
div:hover {
    font-size: 52px; /* Larger font size when hovered */
}
</style>