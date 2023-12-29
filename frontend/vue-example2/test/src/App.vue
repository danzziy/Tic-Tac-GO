<template>
    <div style="z-index: 3;">
        <Board></Board>
        <Button text="Test" color="green"/>
    </div>
</template>

<script>
import Board from './components/Board.vue'
import Button from './components/Button.vue'

export default {
    // reactive state
    name: 'App',
    components: {
        Board,
        Button, // Ensure that 'Board' is added here
    },
}
</script>