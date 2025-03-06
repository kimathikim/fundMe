### how to list contents of collection in the database
```js
const collections = db.getCollectionNames();
collections.forEach(collection => {
  print(`Documents in collection: ${collection}`);
  db[collection].find().forEach(printjson);
  print("\n"); // Add a newline for readability
});
```


# yellow coding robot quiz

Grade 4 Quiz: Yellow Coding & Scratch Basics

Name: ________________________  Date: ______________

Section 1: Multiple Choice Questions

Circle the correct answer.

What is a robot?

a) A type of animal

b) A machine that can be programmed to perform tasks

c) A toy that moves on its own

d) A human assistant

What makes a Yellow Coding Robot special?

a) It can talk to people

b) It can be controlled using LEGO pieces and programmed to do tasks

c) It is only used in factories

d) It does not need electricity

Which component helps the Yellow Coding Robot detect light?

a) Push Button Module

b) Photoresistor Module

c) 7-Color Flashing LED Module

d) PIR Motion Sensor

What does the PIR Motion Sensor do?

a) Makes the robot move faster

b) Detects light changes

c) Detects movement around the robot

d) Controls the robot’s wheels

What is Scratch used for?

a) Cooking food

b) Creating animations and games

c) Controlling a TV

d) Fixing robots

Section 2: True or False

Write T for True or F for False.

___ The 7-Color Flashing LED Module can display different colors based on how the robot is programmed.

___ The Push Button Module is used to start and stop the robot.

___ KidsBlock IDE is a tool used for drawing pictures.

___ The RJ11 Cable is used to connect different parts of the robot together.

___ You can use Scratch to make the Yellow Coding Robot move.

Section 3: Match the Components with Their Functions

Draw a line to match each component to its correct function.

Component

Function

KidsBits Yellow Robot

(a) Detects motion using infrared light

KidsBits Push Button Module

(b) Lights up in different colors

KidsBits 7-Color Flashing LED

(c) Detects light in the environment

KidsBits Photoresistor Module

(d) Main body with motors and wheels

KidsBits PIR Motion Sensor

(e) Starts and stops the robot

Section 4: Label the Images

Write the correct name or function for each image.

[Image of the Yellow Coding Robot]What is the name of this robot?Answer: ___________________________

[Image of Scratch programming blocks]Which software is this?Answer: ___________________________

[Image of a push button module]What is this component called?Answer: ___________________________

[Image of an LED module changing colors]What does this component do?Answer: ___________________________

[Image of a computer connected to the Yellow Coding Robot with a USB cable]Why is the USB cable important?Answer: ___________________________

Bonus Question

What is one cool thing you can make the Yellow Coding Robot do using coding?

Answer: ___________________________________________________

Teacher’s Notes:

This quiz covers basic concepts of robotics, coding, and Scratch.

The image labeling section can be modified based on available images.

Encourage students to be creative in the bonus question.

Answer Key (For Teachers)

b

b

b

c

b

T

T

F

T

T

Matching Answers:

KidsBits Yellow Robot → (d)

KidsBits Push Button Module → (e)

KidsBits 7-Color Flashing LED → (b)

KidsBits Photoresistor Module → (c)

KidsBits PIR Motion Sensor → (a)


