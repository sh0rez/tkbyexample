import React from "react"
import { graphql, Link } from "gatsby"

import { Helmet } from "react-helmet"

import Footer from "../components/footer"

const IndexPage = ({ data }) => (
  <div id="intro">
    <Helmet defer={false}>
      <title>Tanka by Example</title>
      <meta
        name="description"
        value="Hands-on introduction to Grafana Tanka and Jsonnet using annotated
        example programs, inspired by gobyexample.com"
      />
    </Helmet>
    <h2>Tanka by Example</h2>
    <p>
      <a href="https://tanka.dev">Grafana Tanka</a> is a composable
      configuration utility for Kubernetes. It leverages the{" "}
      <a href="https://jsonnet.org">Jsonnet</a> language to realize flexible,
      reusable and concise configuration.
    </p>
    <p>
      <em>Tanka by Example</em> is a hands-on introduction to Tanka and Jsonnet
      using annotated example programs, inspired by the popular{" "}
      <a href="https://gobyexample.com">gobyexample.com</a>. Check out the first
      example or browse the full list below.
    </p>
    <ul>
      {data.allMarkdownRemark.nodes.map(n => (
        <li>
          <Link to={n.frontmatter.path}>{n.frontmatter.title}</Link>
          {n.frontmatter.description ? ": " + n.frontmatter.description : ""}
        </li>
      ))}
    </ul>
    <Footer></Footer>
  </div>
)

export default IndexPage

export const query = graphql`
  {
    allMarkdownRemark {
      nodes {
        frontmatter {
          path
          title
          description
        }
      }
    }
  }
`
