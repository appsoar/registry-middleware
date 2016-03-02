package client

import (
	"fmt"
	//	"strconv"
	"testing"
)

var (
	TestUrl   = "http://192.168.2.110:5000"
	TestImage = "rancher/agent"
	TestTag   = "v0.8.2"
)

func TestCheckVersion(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}
	if err := CheckVersion(opts); err != nil {
		t.Error("CheckVersion 函数调用失败")
		fmt.Println(err)
		//	t.Error(error)
	} else {
		t.Log("测试成功")
	}
}

func TestListRepositories(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}
	n := 2
	if err := ListRepositoriesPagination(opts, n); err != nil {
		t.Log(err)
		t.Error("ListRepositories函数调用失败")
		//	t.Error(error)
	} else {
		t.Log("测试成功")
	}
}

func TestListRepositories2(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}
	n := 0

	if err := ListRepositoriesPagination(opts, n); err != nil {
		t.Log(err)
		t.Error("ListRepositories函数调用失败")
		//	t.Error(error)
	} else {
		t.Log("测试成功")
	}
}

func TestListImageTagsNotExist(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}

	if err := ListImageTags(opts, "xxxxxxxx"); err != nil {
		t.Log("测试成功")
	} else {
		t.Log(err)
		t.Error("TestListImageTags函数调用失败")
	}
}

func TestListImageTagsExist(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}

	if err := ListImageTags(opts, "rancher/agent"); err != nil {
		t.Log(err)
		t.Error("TestListImageTags函数调用失败")
	} else {
		t.Log("测试成功")
	}
}

func TestGetImageManifestsNotExist(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}
	image := "rancher/agent"
	tag := "v0.8"
	if _, err := GetImageManifests(opts, image, tag); err != nil {
		t.Log(err)
		t.Log("测试成功")
	} else {
		t.Error("测试失败")
	}

}

/*
func TestGetImageManifestsExist(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}
	image := "rancher/agent"
	tag := "v0.8.2"
	if manifest, err := GetImageManifests(opts, image, tag); err != nil {
		t.Log(err)
		t.Error("测试失败")
	} else {
		fmt.Println("schemaVersion:" + strconv.Itoa(manifest.SchemaVersion))
		fmt.Println("name:" + manifest.Name)
		fmt.Println("tag:" + manifest.Tag)
		fmt.Println("architecture:" + manifest.Architecture)
		fmt.Println("fslayers:[")
		for i := 0; i < len(manifest.FsLayers); i++ {
			fmt.Println("  blobSum:" + manifest.FsLayers[i].BlobSum)
		}
		fmt.Println("]")
		fmt.Println("history:[")
		for i := 0; i < len(manifest.History); i++ {
			fmt.Println("  ...")
		}
		fmt.Println("]")

		fmt.Println("signatures:[")
		for i := 0; i < len(manifest.Signatures); i++ {
			fmt.Println("  signatures:[")
			fmt.Println("    headers:[")
			fmt.Println("      jwk:[")
			fmt.Println("        crv:" + manifest.Signatures[i].Header.Jwk.Crv)
			fmt.Println("        kid:" + manifest.Signatures[i].Header.Jwk.Kid)
			fmt.Println("        kty:" + manifest.Signatures[i].Header.Jwk.Kty)
			fmt.Println("        x:" + manifest.Signatures[i].Header.Jwk.X)
			fmt.Println("        y:" + manifest.Signatures[i].Header.Jwk.Y)
			fmt.Println("     ]")
			fmt.Println("     alg:" + manifest.Signatures[i].Header.Alg)
			fmt.Println("     ]")
			fmt.Println("   signature:" + manifest.Signatures[i].Signature)
			fmt.Println("   protected:" + manifest.Signatures[i].Protected)
			fmt.Println("  ]")
			fmt.Println("]")
		}

	}

}
*/
func TestGetImageDigestExist(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}
	if digest, err := GetImageDigest(opts, TestImage, TestTag); err != nil {
		t.Error("函数调用失败")
		fmt.Println(err)
		//	t.Error(error)
	} else {
		t.Log(digest)
		t.Log("测试成功")
	}
}

func TestDeleteImage(t *testing.T) {
	opts := ClientOpts{
		Url: TestUrl,
	}
	if err := DeleteImage(opts, TestImage, TestTag); err != nil {
		t.Error("函数调用失败")
		fmt.Println(err)
		//	t.Error(error)
	} else {
		if err := ListRepositoriesPagination(opts, 0); err == nil {
			t.Log("测试成功")
		}
	}
}
