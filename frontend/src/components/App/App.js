import React, { Component, Fragment } from 'react'
import { LandingPage } from '../LandingPage/LandingPage'
import BuildForm from '../Form/Form.js'
import HowItWorks from '../HowItWorks/HowItWorks.js'
import { Route } from 'react-router-dom'
import './App.css'

class App extends Component {
  constructor () {
    super()

    this.state = {}
  }

  render () {

    return (
      <Fragment>
        <Route exact path='/' render={() => (
            <LandingPage />
          )} />
        <Route path='/form' render={() => (
            <BuildForm />
          )} />
         <Route path='/how-it-works' render={() => (
            <HowItWorks />
          )} />
      </Fragment>
    )
  }
}

export default App