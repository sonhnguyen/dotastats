package dotastats

// func GetVideoByLink(site, id, title string, mongodb Mongodb) (Video, error) {
// 	url := URL{Site: site, ID: id}
// 	result, err := GetVideoByLinkMongo(url, mongodb)
// 	if err != nil {
// 		url := URL{Site: site, ID: id}
// 		comment := []Comment{}
// 		result, err = CreateVideoMongo(Video{Url: url, Comment: comment, Title: title}, mongodb)
// 		if err != nil {
// 			return Video{}, err
// 		}
// 	}
// 	return result, nil
// }

// func GetAllVideo(limit, offset string, mongodb Mongodb) ([]Video, error) {
// 	result, err := GetAllVideoMongo(limit, offset, mongodb)
// 	if err != nil {
// 		return []Video{}, err
// 	}
// 	return result, nil
// }

// func GetRandomVideo(limit string, mongodb Mongodb) ([]Video, error) {
// 	result, err := GetRandomVideoMongo(limit, mongodb)
// 	if err != nil {
// 		return []Video{}, err
// 	}
// 	return result, nil
// }

// func GetVideoById(id string, mongodb Mongodb) (Video, error) {
// 	result, err := GetVideoByIdMongo(id, mongodb)
// 	if err != nil {
// 		return Video{}, err
// 	}
// 	return result, nil
// }
// func PostCommentById(id string, comment Comment, mongodb Mongodb) (Video, error) {
// 	result, err := InsertCommentVideoMongo(id, comment, mongodb)
// 	if err != nil {
// 		return Video{}, err
// 	}
// 	return result, nil
// }
