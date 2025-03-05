import React from 'react';
import ReactDOM from 'react-dom/client';
import { BlogLayout } from '../layouts/BlogLayout';
// import sampleContent from './sample-content.js';

// Sample content for development
const DevBlog = () => {
  // Use React Router or a simple state to switch between different sample posts
  // const [currentPost, setCurrentPost] = React.useState(0);
  
  return (
    <div className="dev-container">
      <div className="dev-preview">
        <BlogLayout 
          // title={sampleContent[currentPost].title}
          // content={sampleContent[currentPost].content}
          // currentPath={`/blog/${sampleContent[currentPost].slug}`}
        />
      </div>
    </div>
  );
};

ReactDOM.createRoot(document.getElementById('root')).render(<DevBlog />);
