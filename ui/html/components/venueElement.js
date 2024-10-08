import { LitElement, html, css } from 'lit';
import { Reservation, Venue } from '../models/model.js';
import { getVenues, searchVenues } from '../httpClients/venueAPIClient.js'

export class VenueElement extends LitElement {
  static styles = css`
  
  `;
  static properties = {
    venues: [Venue],
    selectedVenue: { type: Venue }
  };
  constructor() {
    super();
    this.venues = [];
    this.selectedVenue = new Venue;
  }
  connectedCallback() {
    super.connectedCallback();
    getVenues()
      .then(v => this.venues = v)
      .catch(e => {
        console.error('::getAllVenue() error', e);
        this.venues = [];
      });
  }

  searchVenues(queryParam) {
    this.searchVenues(queryParam)
    .then(v => this.venues = v)
      .catch(e => {
        console.error('::searchVenues() error', e);
        this.venues = [];
      });
  }

  render() {
    return html`  
    <div>
      <button @click="${() => this.selectedVenue = new Reservation}">Reset</button>
      <input placeholder="Enter an address, city, or ZIP code" type="text" @change="${(e) => this.searchVenues(e.target.value)}">
      <input placeholder="filter api responses" type="text" @change="${(e) => this.filter(e.target.value)}">
    </div>
    <div><venue-element venue=${JSON.stringify(this.selectedVenue)}/></div>
    <hr>
    <br/>
    <br/>
    <br/>

    ${this.venues.map((it) => {
      return html`      
      <img @click=${() => this.selectedRestObect = it} class="image" src="${it.img}" alt="" />
      <br/>
          <p>${it.id}</p>
          <p>${it.venueId}</p>
      `
    }
    )}
`;
  }

  filter(zipCodeOrName) {
    this.selectedVenue = this.venues
      .find(v => 
        v.zipCode == Number(zipCodeOrName) || 
        v.state.toLocaleLowerCase().match(zipCodeOrName.toLocaleLowerCase()) || 
        v.city.toLocaleLowerCase().match(zipCodeOrName.toLocaleLowerCase())) || 
        new Venue
  }
}
customElements.define('venue-el', VenueElement);
