<template>
    <div class="board">
        <button v-for="(cell, index) in cells" :key="index" @click="handleClick(index)" :textContent="cell">
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
            gameState: '000000000',
        };
    },
    // lifecycle hooks
    mounted() {
        if (WS) {
            WS.onmessage = (event) => {
                console.log("From Board.vue: " + event.data);
                var indexOfColon = event.data.indexOf(':');
                var playerMove = event.data.substring(0, indexOfColon);
                this.gameState = playerMove;
                this.cells.forEach((cell, i) => {
                    if(playerMove[i] != 0){
                        this.cells[i]= playerMove[i];
                        console.log("playermove[i]: " + playerMove[i]);
                    }
                })
                if(event.data.includes('Win') || event.data.includes('Lose')) {
                    var gameOverMessage = `You ${event.data.substring(indexOfColon + 1)}`;
                    console.log("GAME OVER: " + gameOverMessage)
                    this.$emit('game-over', gameOverMessage);
                } 
            };
        } else {
            console.error('Socket not available');
        }
    },
    methods: {
        async handleClick(index) {
            if (this.cells[index] === '') {
                const message = this.gameState.substring(0, index) + this.playerNumber + this.gameState.substring(index + 1);
                WS.send(message);
            }
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
    color: white;
}
</style>
