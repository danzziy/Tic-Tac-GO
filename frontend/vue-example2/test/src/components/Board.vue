<template>
    <div class="board">
        <button v-for="(cell, index) in cells" :key="index" @click="handleClick(index)" :textContent="playerNumber">
            {{ cell }}
        </button>
    </div>
</template>

<script>
import WS from './Websocket.vue'

export default {
    // reactive state
    name: 'Board',
    props: {
        playerNumber: Number,
        socket: null
    },
    data() {
        return {
            cells: Array(9).fill(''), // Represents the cells of the board
        };
    },
    // lifecycle hooks
    mounted() {
        if (WS) {
            WS.onmessage = (event) => {
                console.log("From Board.vue: " + event.data);
            };
        } else {
            console.error('Socket not available');
        }
    },
    methods: {
        handleClick(index) {
            console.log(index)
            console.log("PlayerNumber: " + this.playerNumber)
            WS.send("000010000")
        },
    },
}
</script>

<style scope>
.board {
    display: grid;
    grid-template-columns: repeat(3, 100px);
    grid-gap: 0px;
    position: relative;
    pointer-events: all;
}

.board button {
    width: 100px;
    height: 100px;
    background: transparent;
    border: 1px solid white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 2em;
    cursor: pointer;
}
</style>
