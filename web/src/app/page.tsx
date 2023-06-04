'use client'

import { useEffect, useState } from "react";

const Home = () => {
  const [gameState, setGameState] = useState({
    Bomb: "",
    Map: "",
  })

  useEffect(() => {
    let tmp = "";
    switch(gameState.Bomb) {
      case "planted":
        tmp = "#0F0";
        break;
      case "exploded":
        tmp = "#F00";
        break;
      default:
        tmp = "#000";
    }
    document.documentElement.style.setProperty( `--bomb-state`, tmp)
  }, [gameState])

  useEffect(() => {
    let socket = new WebSocket("ws://127.0.0.1:8080/ws");
    console.log("Attempting Connection...");

    socket.onopen = () => {
      console.log("Successfully Connected");
    };

    socket.onclose = event => {
      console.log("Socket Closed Connection: ", event);
      socket.send("Client Closed!")
    };

    socket.onerror = error => {
      console.log("Socket Error: ", error);
    };

    socket.onmessage = ({data}) => {
      const state = JSON.parse(data)
      console.log(state)
      setGameState(prev => ({...prev, ...state}))
    }

    return () => socket.close()

  }, [])

  return (
    <>
    <h1>{gameState.Bomb != "" ? gameState.Bomb : "Bomb is not planted!"}</h1>
    <h1>{gameState.Map}</h1>
    </>
  )
}

export default Home;
