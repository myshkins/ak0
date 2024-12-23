import "./ProjectCard.css";

const ProjectCard = (props) => {
  return (
    <div className="project-card">
      <a href={props.link}>
        <h2>{props.title}</h2>
      </a>
      <div className="project-card-inner">
        <a href={props.link}>
          <img className="project-card-img" src={props.img} />
        </a>
        {props.text}
      </div>
    </div>
  );
};

export default ProjectCard;
