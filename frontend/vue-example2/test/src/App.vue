<template>
    <!-- <div>
        <Board></Board>
        <Button text="Test" color="green"/>
    </div> -->
    <div>
        <div id="home_page" v-if="page=='home_page'">
            <div class="centered-div">TIC-TAC-GO</div>            
            <Button @click="startGame" text="Start Game" />
        </div>

        <div v-else-if="page=='Waiting for Player'">
            <div class="centered-div">Loading...</div>
        </div>

        <div v-else-if="page=='Start Game'">
            <Board />
        </div>
  </div>
</template>

<script>
import Board from './components/Board.vue'
import Button from './components/Button.vue'
import { ref, provide } from 'vue';

export default {
    // reactive state
    name: 'App',
    components: {
        Board,
        Button, // Ensure that 'Board' is added here
    },
    
    setup() {
        let webSocket = null;
        let page=ref('')
        page.value = 'home_page'
        const startGame = async () => {
            webSocket = new WebSocket('ws://localhost:8081/public');
            provide('webSocket', webSocket);

            webSocket.onopen = (event) => {
                console.log('WebSocket connection established.');

                // You can perform actions here upon successful connection
                // For example, sending a 'Join Room' message
                webSocket.send('Join Room');
            };

            webSocket.onmessage = (event) => {
                page.value = event.data;
                console.log("From app.vue: " + page.value);
            }
            webSocket.onerror = (error) => {
                console.error('WebSocket error:', error);
            };
            webSocket.onclose = () => {
                console.log('WebSocket connection closed');
                connected.value = false;
            };
        };
        
        return{startGame, page};
    },
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