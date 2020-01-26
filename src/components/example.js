import React from "react"
import { graphql } from "gatsby"
import rehypeReact from "rehype-react"
import * as clipboard from "clipboard-polyfill"

import "../site.css"
import Footer from "./footer"

const CopyButton = props => (
  <button
    onClick={() => {
      clipboard
        .writeText(atob(props.code))
        .then(() => {
          console.log("copied!")
        })
        .catch(err => {
          console.error(err)
        })
    }}
  >
    {props.children}
  </button>
)

const renderAst = new rehypeReact({
  createElement: React.createElement,
  components: {
    "copy-button": CopyButton,
  },
}).Compiler

const Template = ({ data }) => {
  const { markdownRemark } = data // data.markdownRemark holds your post data
  const { htmlAst, frontmatter } = markdownRemark
  return (
    <div className="center">
      <main>
        {renderAst(htmlAst)}

        <div>
          <Footer></Footer>
        </div>
      </main>
    </div>
  )
}

export default Template

export const pageQuery = graphql`
  query($path: String!) {
    markdownRemark(frontmatter: { path: { eq: $path } }) {
      htmlAst
      frontmatter {
        path
        title
      }
    }
  }
`
