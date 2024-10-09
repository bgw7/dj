import { LitElement, html, css } from 'lit';
import { Track } from '../models/model.js';
import { getTracks, createTrack, createVote, deleteVote } from '../httpClients/djroombaAPIClient.js'

export class DJRoombaElement extends LitElement {
  static styles = css`
  
  `;
  static properties = {
    songTracks: [Track],
    createTrackErr: Object,
  };
  constructor() {
    super();
    this.songTracks = [];
    this.createTrackErr = null;
    this.createVoteErr = null;
    this.deleteVoteErr = null;
  }
  connectedCallback() {
    super.connectedCallback();
    getTracks()
      .then(v => this.songTracks = v)
      .catch(e => {
        console.error('::getTracks() error', e);
        this.songTracks = [];
      });
  }

  createTrack(trackObj) {
    this.createTrack(trackObj)
      .then(console.log)
      .catch(e => {
        console.error('::createTrack() error', e);
        this.createTrackErr = e;
      });
  }

  createVote(voteObj) {
    this.createVote(voteObj)
      .then(console.log)
      .catch(e => {
        console.error('::createVote error', e);
        this.createVoteErr = e;
      })
  }
  deleteVote(voteObj) {
    this.createVote(voteObj)
      .then(console.log)
      .catch(e => {
        console.error('::createVote error', e);
        this.deleteVoteErr = e;
      })
  }

  render() {
    return html`  
    <div>
    <hr>
    <br/>

    ${this.songTracks.map((it) => {
      return html`      
      <button @click=${() => this.createVote(it.id)}>vote</button>
      <br/>
          <p>${it.id}</p>
          <p>${it.url}</p>
          <p>${it.filename}</p>
          <p>${it.voteCount}</p>
          <p>${it.hasPlayed}</p>
          <p>${it.createdBy}</p>
      `
    }
    )}
`;
  }
}
customElements.define('dj-el', DJRoombaElement);
