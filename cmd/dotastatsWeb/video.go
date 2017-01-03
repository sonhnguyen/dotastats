package main

// func (a *App) GetVideoByLinkHandler() HandlerWithError {
// 	return func(w http.ResponseWriter, req *http.Request) error {

// 		queryValues := req.URL.Query()
// 		site := queryValues.Get("site")
// 		id := queryValues.Get("id")
// 		title := queryValues.Get("title")

// 		video, err := youtime.GetVideoByLink(site, id, title, a.mongodb)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}

// 		err = json.NewEncoder(w).Encode(video)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}
// 		return nil
// 	}
// }

// func (a *App) GetAllVideoHandler() HandlerWithError {
// 	return func(w http.ResponseWriter, req *http.Request) error {

// 		queryValues := req.URL.Query()
// 		limit := queryValues.Get("limit")
// 		offset := queryValues.Get("offset")

// 		video, err := youtime.GetAllVideo(limit, offset, a.mongodb)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}

// 		err = json.NewEncoder(w).Encode(video)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}
// 		return nil
// 	}
// }

// func (a *App) GetRandomVideoHandler() HandlerWithError {
// 	return func(w http.ResponseWriter, req *http.Request) error {

// 		queryValues := req.URL.Query()
// 		limit := queryValues.Get("limit")

// 		video, err := youtime.GetRandomVideo(limit, a.mongodb)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}

// 		err = json.NewEncoder(w).Encode(video)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}
// 		return nil
// 	}
// }

// func (a *App) GetVideoByIdHandler() HandlerWithError {
// 	return func(w http.ResponseWriter, req *http.Request) error {
// 		params := GetParamsObj(req)
// 		id := params.ByName("id")

// 		video, err := youtime.GetVideoById(id, a.mongodb)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}

// 		err = json.NewEncoder(w).Encode(video)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}

// 		return nil
// 	}
// }

// func (a *App) PostCommentByIdHandler() HandlerWithError {
// 	return func(w http.ResponseWriter, req *http.Request) error {
// 		var comment youtime.Comment

// 		err := json.NewDecoder(req.Body).Decode(&comment)
// 		if err != nil {
// 			a.logr.Log("error decode param: %s", err)
// 			return newAPIError(400, "error param: %s", err)
// 		}

// 		params := GetParamsObj(req)
// 		id := params.ByName("id")

// 		video, err := youtime.PostCommentById(id, comment, a.mongodb)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}

// 		err = json.NewEncoder(w).Encode(video)
// 		if err != nil {
// 			a.logr.Log("error when return json %s", err)
// 			return newAPIError(404, "error when return json %s", err)
// 		}

// 		return nil
// 	}
// }
