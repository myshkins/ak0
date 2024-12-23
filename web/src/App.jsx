import "./App.css";
import blurbs from './components/Blurbs'
import NavBar from "./components/NavBar";
import Hero from "./components/Hero";
import SkillsBox from "./components/SkillsBox";
import Projects from "./components/Projects";
import ProjectCard from "./components/ProjectCard";

function App() {
  return (
    <div className="app">
      <NavBar />
      <Hero>
        <SkillsBox />
      </Hero>
      <Projects>
        <ProjectCard
          link="https://jupiter.ak0.io"
          title={"jupiter coffee"}
          img="https://ak0-io.us-east-1.linodeobjects.com/screenshot-jupiter.webp"
          text={blurbs.jupiter} />
        <ProjectCard
          link="https://airboo.ak0.io"
          title={"airboo"}
          img="https://ak0-io.us-east-1.linodeobjects.com/screenshot-airboo.webp"
          text={blurbs.airboo} />
        <ProjectCard
          link="https://gopetwatch.ak0.io"
          title={"pet-watch"}
          img="https://ak0-io.us-east-1.linodeobjects.com/screenshot-petwatch.webp"
          text={blurbs.petwatch} />
        <ProjectCard
          link="https://blog.ak0.io"
          title={"my blog"}
          img="https://ak0-io.us-east-1.linodeobjects.com/screenshot-blog.webp"
          text={blurbs.blog} />
        <ProjectCard
          link="https://etch.ak0.io"
          title={"etch-a-sketch"}
          img="https://ak0-io.us-east-1.linodeobjects.com/screenshot-etch.webp"
          text={blurbs.etch} />
      </Projects>
    </div>
  );
}

export default App;
