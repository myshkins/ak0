// scripts/prerender-components.jsx
console.log('Current working directory:', process.cwd());

import React from 'react';
import ReactDOMServer from 'react-dom/server';
import fs from 'fs';
import path from 'path';
import { BlogLayout } from '../web/src/layouts/BlogLayout'



// Create the output directory
const outputDir = path.resolve(__dirname, '../internal/handlers/dist/posts/');
if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}


// Render blog layout (without content)
const blogLayoutHtml = ReactDOMServer.renderToStaticMarkup(
  <BlogLayout >
    {/* {{.Content}} */}
  </BlogLayout>
);

fs.writeFileSync(
  `${outputDir}/blog-layout.html`,
  blogLayoutHtml
);

console.log('Prerendering complete!');
