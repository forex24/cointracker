import {
    GET_LIST,
    GET_ONE,
    CREATE,
    UPDATE,
    DELETE,
    DELETE_MANY,
    GET_MANY,
    GET_MANY_REFERENCE,
  } from 'react-admin';
  
  import { stringify } from 'qs';
  import merge from 'deepmerge';
  import axios from 'axios';
  
  import defaultSettings from './default-settings';
  import { NotImplementedError } from './errors';
  import init from './initializer';
  
  // Set HTTP interceptors.
  init();
  
  /**
   * Maps react-admin queries to a JSONAPI REST API
   *
   * @param {string} apiUrl the base URL for the JSONAPI
   * @param {string} userSettings Settings to configure this client.
   *
   * @param {string} type Request type, e.g GET_LIST
   * @param {string} resource Resource name, e.g. "posts"
   * @param {Object} payload Request parameters. Depends on the request type
   * @returns {Promise} the Promise for a data response
   */
  export default (apiUrl, userSettings = {}) => (type, resource, params) => {
    let url = 'api/v1';
    const settings = merge(defaultSettings, userSettings);
  
    const options = {
      headers: settings.headers,
    };
  
    switch (type) {
        case GET_LIST: {
            const { page, perPage } = params.pagination;
            const query = {
                'offset': (page - 1) * perPage,
                'limit': perPage,
                'text': params.filter.q,
            };
            url = `${apiUrl}/${resource}?${stringify(query)}`;
            break;
        }
  
        case GET_ONE:
            url = `${apiUrl}/${resource}/${params.id}`;
            break;

        case CREATE:
            url = `${apiUrl}/${resource}`;
            options.method = 'POST';
            var data = params.data

            if (data['file'] !== undefined) {
              let form = new FormData()
              var file = data.file
              if (file && file.rawFile) {
                form.append('file', file.rawFile)
              }
              options.headers = {'Content-Type': 'multipart/form-data'}
              options.data = form
              url += "/import"
            } else {
              options.data = JSON.stringify(params.data);
            }
            break;

        case UPDATE: {
            url = `${apiUrl}/${resource}/${params.id}`;


            options.method = 'PUT';
            options.data = JSON.stringify(params.data);
            break;
        }
        case DELETE_MANY:
          url = `${apiUrl}/${resource}`;
          options.data = JSON.stringify(params.ids);
          options.method = 'DELETE';
          break;
        case DELETE:
          url = `${apiUrl}/${resource}/${params.id}`;
          options.method = 'DELETE';
          break;
        case GET_MANY:
          var query = {
              'ids': params.ids,
          };
          url = `${apiUrl}/${resource}?${stringify(query)}`;
          break;
          
        case GET_MANY_REFERENCE:
          var query = {}
          query[params.target] = params.id
          url = `${apiUrl}/${resource}?${stringify(query)}`;
          break;
        default:
        throw new NotImplementedError(`Unsupported Data Provider request type ${type}`);
    }
  
    return axios({ url, ...options })
      .then((response) => {
        switch (type) {
          case GET_LIST: {
            return {
              data: response.data.data,
              total: response.data.meta[settings.total],
            };
          }
  
          case GET_ONE: {
            return {
              data: response.data.data,
              meta: response.data.meta,
            };
          }
  
          case GET_MANY: {
            return {
              data: response.data.data,
              total: response.data.meta[settings.total],
            };
          }

          case GET_MANY_REFERENCE: {
            return {
              data: response.data.data,
              total: response.data.meta[settings.total],
            };
          }
  
          case CREATE: {
            
            const { id, attributes } = response.data.data;
            return {
              data: {
                id, ...attributes,
              },
            };
          }
  
          case UPDATE: {
            const { id, attributes } = response.data.data;
  
            return {
              data: {
                id, ...attributes,
              },
            };
          }
          
          case DELETE_MANY:
            return {
              data: response.data.data,
              meta: response.data.meta,
            };
          case DELETE: {
            return {
              data: "DELETE"
            };
          }
  
          default:
            throw new NotImplementedError(`Unsupported Data Provider request type ${type}`);
        }
      });
  };