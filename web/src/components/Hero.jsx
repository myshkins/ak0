import "./Hero.css";
import Selfie from "../assets/pixel-selfie.png";
const Hero = (props) => {
  return (
    <div className="hero">
      <div className="selfie-wrap">
        <img className="selfie" src={Selfie} />
      </div>
      <div className="hero-desktop-text">
        Hey, I'm Alex Krenitsky, a NYC public school teacher turned software
        developer. These are the things I've built. I'm always eager to connect
        , so let me know if you like anything that you see!
      </div>
      {props.children}
    </div>
  );
};

export default Hero;
