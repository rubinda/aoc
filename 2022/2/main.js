#!/usr/bin/node
const fs = require("fs");

const objectCodes = new Map([
  ["A", "rock"],
  ["B", "paper"],
  ["C", "scissors"],
  // Part 1 of the challenges gives you the object played
  ["X", "rock"],
  ["Y", "paper"],
  ["Z", "scissors"],
]);
// Part 2 of the challenge gives you the desired outcome
const outcomeCodes = {
  X: "loss",
  Y: "draw",
  Z: "win",
};
const objectValues = new Map([
  ["rock", 1],
  ["paper", 2],
  ["scissors", 3],
]);
const obArray = Array.from(objectValues.keys());
const outcomeValues = {
  win: 6,
  draw: 3,
  loss: 0,
};

const playObject = (move) => {
  let [a, _, b] = move.split("");
  a = objectValues.get(objectCodes.get(a));
  b = outcomeCodes[b];
  return b === "win"
    ? obArray.at(a % 3)
    : b === "loss"
    ? obArray.at((a - 2) % 3)
    : obArray.at(a - 1);
};
const moveOutcome = (move) => {
  let [a, _, b] = move.split("");
  a = objectCodes.get(a);
  b = objectCodes.get(b);
  let w = objectValues.get(a) % 3;
  return objectValues.get(a) - objectValues.get(b) === 0
    ? outcomeValues.draw
    : b === obArray.at(w)
    ? outcomeValues.win
    : outcomeValues.loss;
};

fs.readFile("example.in", "utf8", (err, data) => {
  if (err) {
    console.error(err);
    return;
  }
  const c1 = data.split("\n").reduce((prev, curr) => {
    return (
      prev + objectValues.get(objectCodes.get(curr[2])) + moveOutcome(curr)
    );
  }, 0);

  const c2 = data.split("\n").reduce((prev, curr) => {
    return (
      prev +
      outcomeValues[outcomeCodes[curr[2]]] +
      objectValues.get(playObject(curr))
    );
  }, 0);

  console.log(c1, c2);
});
