const blurbs = {
  airboo: [
    `My brother and I wanted to collaborate on a project, and us both being 
    environmentalists at heart, we decided to build an air quality monitoring 
    app. The end goal is to create a platform where users can compare air 
    quality measurements they personally took with air quality data provided by 
    government agencies. We first started with building the data pipeline and 
    data visualization tool that you see linked here. We used Apache 
    Airflow to orchestrate our ELT data pipeline. Our api is built with FastAPI, 
    and the front end is rendered with React. Click `, 
    <a href="https://airboo.ak0.io">here</a>,
    ` to see the app, or the `,
    <a href="https://github.com/myshkins/airboo">Github repo</a>,
    ` for this project.`
  ],
  blog: [
    `I built a blog for recording all of the things I learn as I develop 
    projects. It's build with Flask, Postgres, and SQLAlchemy, and it's deployed 
    with Nginx, Gunicorn, and Docker. Check out the site `,
    <a href="https://blog.ak0.io">here</a>,
    " or see the ",
    <a href="https://github.com/myshkins/ak0_blog">Github repo</a>, "."
  ],
  etch: [
    `Here is a fun little drawing app created with JavaScript, CSS, and HTML. 
    In this app a user only needs to move their mouse to draw on the screen. 
    Users can also choose the dimensions of the grid (ie. the resolution), as 
    well as the color of the pen. Clicking the rainbow button allows the pen or 
    the grid to be drawn in a rainbow of colors. Click `,
    <a href="https://etch.ak0.io">here</a>,
    " to see the app or the ",
    <a href="https://github.com/myshkins/etch_sketch">Github repo</a>,
    "."
  ],
  jupiter: [
    `Work In Progress! An exercise in all things front end. I wanted to show the
    world that I too can center a div. jk. This was more about learning how to
    build an SPA that's responsive, server side generated, and (at least somewhat)
    nice looking. I based this fictional coffee shop website loosely on Blue Bottle
    Coffee's website. Next steps are to add more pages and a guest user login. But,
    hey, in the meantime, stop on by this cute `,
    <a href="https://jupiter.ak0.io">lil internet cafe</a>,
    ` and grab yourself a fancy $20 coffee. It'll be a nice experience for you on
    any device. And here's the `,
    <a href="https://github.com/myshkins/koffee2">Github repo</a>,
    `.`,
  ],
  petwatch: [
    `*UPDATED: This project was previously written in Python, but as an exercise
    I rewrote it in Go! Checkout the new version :)
    The inspiration for this project came from Chancho the chinchilla. He's a 
    15 year old chinchilla who is very particular about the climate around him. 
    Chinchillas generally don't tolerate temperatures outside the range of 60-75 
    degrees Farenheit. Thus, I created this app as a way to monitor the 
    temperature of the apartment when my partner and I are away. I used a 
    Raspberry Pi with a temperature sensor to get the temperature data, which 
    is then posted to my database via my api. The app then pulls the data to 
    render a webpage with a view of the data. The tech stack includes Go, Gin, 
    MYSQL, Docker, and Nginx. It's deployed on my personal VPS. Click `,
    <a href="https://gopetwatch.ak0.io">here</a>,
    ` to see the app, or the `,
    <a href="https://github.com/myshkins/gopetwatch">Github repo</a>,
    ` for this project.`
  ],
}

export default blurbs
