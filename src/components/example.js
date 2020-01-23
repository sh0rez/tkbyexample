import React from "react"
import { graphql } from "gatsby"

import "../site.css"

const Template = ({ data }) => {
  const { markdownRemark } = data // data.markdownRemark holds your post data
  const { html, frontmatter } = markdownRemark
  return (
    <div className="center">
      <main>
        <div dangerouslySetInnerHTML={{ __html: html }} />

        <div>
          <p className="footer">
            Inspired by <a href="https://gobyexample.com">gobyexample.com</a> |{" "}
            <a href="https://github.com/sh0rez/tkbyexample">source</a> |{" "}
            <a href="https://github.com/sh0rez/tkbyexample/blob/master/LICENSE">
              license
            </a>
          </p>
        </div>
      </main>
    </div>
  )
}

export default Template

export const pageQuery = graphql`
  query($path: String!) {
    markdownRemark(frontmatter: { path: { eq: $path } }) {
      html
      frontmatter {
        path
        title
      }
    }
  }
`
