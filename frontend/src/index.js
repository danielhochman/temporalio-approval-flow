import React from "react";
import ReactDOM from 'react-dom';
import {useForm} from "react-hook-form";
import {CheckCircle, ChevronDown, ChevronRight, Lock, MessageCircle, Unlock} from "react-feather";

const StatusIcon = ({status}) => {
    switch (status) {
        case "lock":
            return <Lock color={"#c62828"}/>
        case "unlock":
            return <Unlock color={"olive"}/>
        case "comment":
            return <MessageCircle color={"gray"}/>
        case "approve":
            return <CheckCircle color={"green"}/>
        default: {
            return "?"
        }
    }
}

const WorkflowRow = ({workflow, setStale}) => {
    const [expanded, setExpanded] = React.useState(false);
    const {register, handleSubmit, formState, reset} = useForm();
    s
    React.useEffect(() => {
        if (formState.isSubmitSuccessful) {
            reset();
            setStale(true);
        }
    }, [formState]);

    const onSubmit = async v =>
        await fetch(
            "/api/AddReview",
            {method: "POST", body: JSON.stringify(v)}).then(
            () => {
            }, err => console.log(err))

    let Status = () => <ins>Open</ins>
    if (workflow["status"] === "locked") {
        Status = () => <span style={{color: "#c62828"}}><b>Locked</b></span>
    } else if (workflow["status"] === "approved") {
        Status = () => <i>
            <ins>Approved</ins>
        </i>
    }

    return <>
        <tr onClick={() => setExpanded(!expanded)}>
            <td>{expanded ? <ChevronDown/> : <ChevronRight/>} </td>
            <td>{workflow["startTime"]}</td>
            <td>{workflow["user"]}</td>
            <td>{workflow["action"]}</td>
            <td>
                <Status/>
            </td>
            <td style={{textAlign: "right"}}>{workflow["reviews"].length}</td>
        </tr>
        {expanded && <tr>
            <td/>
            <td colSpan={5}>
                <div><strong><small>Details</small></strong></div>
                <div>
                    <table>
                        <tbody>
                        <tr>
                            <th scope={"row"} width={"16%"}>Workflow ID</th>
                            <td>{workflow["id"]}</td>
                        </tr>
                        <tr>
                            <th scope={"row"}>Run ID</th>
                            <td>{workflow["runId"]}</td>
                        </tr>
                        </tbody>
                    </table>
                </div>
                <div><strong><small>Comments</small></strong></div>
                <div>
                    <table>
                        <tbody>
                        {workflow["reviews"].length === 0 && <tr>
                            <td><i>No comments yet.</i></td>
                        </tr>}
                        {workflow["reviews"].map(r => <tr>
                            <td width={"16%"}>{r["time"]}</td>
                            <td width={"10%"}>{r["user"]}</td>
                            <td width={"5%"}><StatusIcon status={r["status"]}/></td>
                            <td>{r["message"]}</td>
                        </tr>)}
                        <tr>
                            <td colspan={4}><br/> <strong>Submit a comment</strong></td>
                        </tr>
                        <tr>
                            <td colspan={4}>
                                <form onSubmit={handleSubmit(onSubmit)}>
                                    <input type="hidden" {...register("id", {value: workflow["id"]})} />
                                    <input type="hidden" {...register("runId", {value: workflow["runId"]})} />
                                    <textarea style={{resize: "none"}}
                                              placeholder={"Enter comment text here"} {...register("message", {required: true})} />
                                    <fieldset>
                                        <legend><b>Action</b></legend>
                                        <label>
                                            <input {...register("action", {required: true})} type="radio"
                                                   value="comment"/>
                                            <b><StatusIcon status={"comment"}/> Comment</b>
                                        </label>
                                        <label>
                                            <input {...register("action", {required: true})} type="radio"
                                                   value="approve"/>
                                            <b><StatusIcon status={"approve"}/> Approve</b>
                                        </label>
                                        <label>
                                            <input {...register("action", {required: true})} type="radio"
                                                   value="lock"/>
                                            <b><StatusIcon status={"lock"}/> Lock</b>
                                        </label>
                                        <label>
                                            <input {...register("action", {required: true})} type="radio"
                                                   value="unlock"/>
                                            <b><StatusIcon status={"unlock"}/> Unlock</b>
                                        </label>
                                    </fieldset>
                                    <button className={"outline"} type={"submit"}>Submit comment</button>
                                </form>
                            </td>
                        </tr>
                        </tbody>
                    </table>
                </div>
            </td>
        </tr>}
    </>
}

const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const App = () => {
    const [workflows, setWorkflows] = React.useState([]);
    const [submitting, setSubmitting] = React.useState(false);
    const [loading, setLoading] = React.useState(false);
    const [stale, setStale] = React.useState(false);
    React.useEffect(() => {
        setLoading(true)
        fetch("/api/ListOpenWorkflow").then(r => r.json()).then(r => {
            setWorkflows(r["workflows"])
        }, err => {
            console.log(err)
        }).finally(() => setLoading(false))
    }, [])

    const {register, handleSubmit, formState, reset} = useForm();

    React.useEffect(() => {
        if (formState.isSubmitSuccessful) {
            reset();

        }
        if (formState.isSubmitSuccessful || stale) {
            setLoading(true)
            fetch("/api/ListOpenWorkflow").then(r => r.json()).then(r => {
                setWorkflows(r["workflows"])
            }, err => {
                console.log(err)
            }).finally(() => {
                setLoading(false);
                setStale(false)
            })
        }
    }, [formState, stale]);

    const onSubmit = async v => {
        setSubmitting(true)
        await fetch(
            "/api/ExecuteWorkflow",
            {method: "POST", body: JSON.stringify(v)}).then(
            () => {
            }, err => console.log(err)).finally(() => {

        })
        await sleep(3000) // HACK HACK HACK! Can't see workflow after executing it immediately.
        setSubmitting(false)

    }

    return <main className="container">
        <div>&nbsp;</div>
        <hgroup>
            <h2>Temporal Approval Management System</h2>
            <h3>This is a demo of Temporal built in Go and React by <a href={"https://github.com/danielhochman/temporalio-approval-flow"}>@danielhochman</a>.</h3>
        </hgroup>
        <form onSubmit={handleSubmit(onSubmit)}>
            <div className="grid">
                <label>Username <input autoComplete="off" disabled={formState.isSubmitting} {...register("user")} /></label>
                <label>Action
                    <select {...register("action")} required disabled={formState.isSubmitting}>
                        <option value="" selected>Select an action…</option>
                        <option>Terminate an instance</option>
                        <option>Create a repo</option>
                        <option>Resize an ASG</option>
                    </select>
                </label>
            </div>
            <button aria-busy={submitting} type="submit">Create Approval Request</button>
        </form>
        <h4>Open Workflows ({workflows.length}) <span aria-busy={loading}/></h4>
        <table>
            <thead>
            <tr>
                <th scope="col" width={"1%"}/>
                <th scope="col" width={"25%"}>Time (UTC)</th>
                <th scope="col">User</th>
                <th scope="col">Action</th>
                <th scope="col">Status</th>
                <th scope="col" style={{textAlign: "right"}}>Comments</th>
            </tr>
            </thead>
            <tbody>
            {workflows.map(w => <WorkflowRow key={w["id"] + "_" + w["runId"]} workflow={w} setStale={setStale}/>)}
            </tbody>
        </table>
    </main>
}

const app = document.getElementById("app");
ReactDOM.render(<App/>, app);
