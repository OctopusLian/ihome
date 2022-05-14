/*
 * @Description:
 * @Author: neozhang
 * @Date: 2022-05-14 16:36:35
 * @LastEditors: neozhang
 * @LastEditTime: 2022-05-14 16:37:43
 */
package utils

func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url
	return domain_url
}
