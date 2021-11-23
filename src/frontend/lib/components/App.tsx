import React from 'react'
import { Switch, Route } from 'react-router-dom'

const Deposit = () => (
    <div>
        <form encType="multipart/form-data" method="POST" action="/deposit">
            <label>File: </label><input name="uploadfile" type="file" value="" />
            <br />
            <input type="submit" value="Deposit" />
        </form>
    </div>
)

const App = () => (
    <Switch>
        <Route exact path={'/app/deposit'}><Deposit /></Route>
        <Route exact path={'/app'} render={() => (<div><h1>Yes, this is frontend.</h1></div>)} />
        <Route path={'*'} render={() => (<div><h1>404.</h1></div>)} />
    </Switch>
)

export default App