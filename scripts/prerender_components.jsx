// scripts/prerender-components.jsx
import React from 'react';
import ReactDOMServer from 'react-dom/server';
import fs from 'fs';
import path from 'path';
import { BlogLayout } from '../web/src/components/BlogLayout';

// Create the output directory
const outputDir = path.resolve(__dirname, '../internal/handlers/templates/components');
if (!fs.existsSync(outputDir)) {
  fs.mkdirSync(outputDir, { recursive: true });
}


// Render blog layout (without content)
const blogLayoutHtml = ReactDOMServer.renderToStaticMarkup(
  <BlogLayout title="{{.Title}}" currentPath="{{.CurrentPath}}">
    {/* {{.Content}} */}
  </BlogLayout>
);

fs.writeFileSync(
  `${outputDir}/blog-layout.html`,
  blogLayoutHtml
);

console.log('Prerendering complete!');
