import React from 'react'
import { Switch, Route, BrowserRouter as Router } from 'react-router-dom'

const Deposit = () => (
    <div>
        <form encType="multipart/form-data" method="POST" action="/api/deposit">
            <label>File: </label><input name="uploadfile" type="file" />
            <br />
            <input type="submit" value="Deposit" />
        </form>
    </div>
)

const App = () => (
    <Router>
        <Switch>
            <Route exact path={'/app/deposit'}><Deposit /></Route>
            <Route exact path={'/app'} render={() => (<div><h1>Yes, this is frontend.</h1></div>)} />
            <Route path={'*'} render={() => (<div><h1>404.</h1></div>)} />
        </Switch>
    </Router>
)

export default App