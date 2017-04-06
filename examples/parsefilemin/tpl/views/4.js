


function hello_world(name) {
    // Check to see if name is valid, otherwise were are going to say
    // Hello to the world.
    if ( name ) {
        // It's always good to get to known people more personally.
        console.log("hello " + name + "!");
    } else {
        // It's nice when we are talking to the world; we can be generic.
        console.log("Hello World!")
    }

    /*
    These comments should be stripped out as they only add bulk to this javascript file.
     */

}

// Let's say Hello to the world
hello_world()

// Now let's say hello to gautam.
hello_world("Gautam")



