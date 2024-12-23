import "./Projects.css";

const Projects = (props) => {
  return (
    <>
      <div className="hr"></div>
      <div className="projects">{props.children}</div>{" "}
    </>
  );
};

export default Projects;
