<template>
    <div>
        <div id="home_page" v-if="page === 'Home Page'"> 
            <div class="centered-div">TIC-TAC-GO</div>            
            <Button @click="startGame" text="Start Game" />
        </div>

        <div v-else-if="page === 'Waiting for Player'">
            <div class="centered-div">Loading...</div>
        </div>

        <div v-else-if="page === 'Start Game'">
            <Board  @game-over="gameOver" :playerNumber="playerNumber"/>
        </div>

        <div v-else-if="page === 'Game Over'">
            <div class="centered-div">{{ playerOutcome }}</div>
            <Button @click="startGame" text="End Game" />
        </div>
  </div>
</template>

<script>
import Board from './components/Board.vue'
import Button from './components/Button.vue'
import WS from './components/Websocket.vue'

export default {
    // reactive state
    name: 'App',
    components: {
        Board,
        Button,
    },
    data() {
        return{
            page:  "Home Page",
            playerNumber: 0,
            playerOutcome: '',
        }
    },
    mounted() {
        WS.onopen = (event) => {
            console.log('WebSocket connection established.');
        };

        WS.onmessage = (event) => {
            console.log("APP VUE")
            this.page = event.data;
            
            console.log("From app.vue: " + this.page);
            if( this.playerNumber == 0 && this.page == 'Start Game'){
                this.playerNumber = 2
            } else {
                this.playerNumber = 1
            }
        }
        WS.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        WS.onclose = () => {
            console.log('WebSocket connection closed');
            // connected.value = false;
        };
    },
    methods: {
        async startGame() {
           WS.send("Join Room")
        },
        gameOver(gameOverMessage) {
            this.page = 'Game Over';
            this.playerOutcome = gameOverMessage;
        }
    }
}
</script>

<style>
div {
        position: relative;
}
@font-face {
    font-family: '8bit'; /* Font family name */
    src: url("../public/1up.ttf") /* URL to the font file */
}

.centered-div {
    padding: 20px;
    border-radius: 5px;
    font-family: '8bit';
    font-size: 150px;
    color: white;
    position: fixed;

    justify-content: center;
    align-items: center;
    transform: translate(30%, 40%);
    pointer-events: none;

}
</style>