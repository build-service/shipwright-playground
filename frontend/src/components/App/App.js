import React, { Component, Fragment } from 'react'
import BuildForm from '../Form/Form.js'
import { Route } from 'react-router-dom'
import './App.css'

import { LandingPage } from '../LandingPage/LandingPage'


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
      </Fragment>
    )
  }
}

export default App