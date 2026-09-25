import{Ci as e,Cr as t,ca as n,kr as r,la as i,sr as a,ta as o,ur as s,wi as c,wr as l}from"./router-CZSFYymG.js";import{n as u}from"./index-C1zbMUbF.js";var d=e(`statistic`,[c(`label`,`
 font-weight: var(--n-label-font-weight);
 transition: .3s color var(--n-bezier);
 font-size: var(--n-label-font-size);
 color: var(--n-label-text-color);
 `),e(`statistic-value`,`
 margin-top: 4px;
 font-weight: var(--n-value-font-weight);
 `,[c(`prefix`,`
 margin: 0 4px 0 0;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-prefix-text-color);
 `,[e(`icon`,{verticalAlign:`-0.125em`})]),c(`content`,`
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-text-color);
 `),c(`suffix`,`
 margin: 0 0 0 4px;
 font-size: var(--n-value-font-size);
 transition: .3s color var(--n-bezier);
 color: var(--n-value-suffix-text-color);
 `,[e(`icon`,{verticalAlign:`-0.125em`})])])]),f=n({name:`Statistic`,props:Object.assign(Object.assign({},a.props),{tabularNums:Boolean,label:String,value:[String,Number]}),slots:Object,setup(e){let{mergedClsPrefixRef:n,inlineThemeDisabled:r,mergedRtlRef:i}=l(e),c=a(`Statistic`,`-statistic`,d,u,e,n),f=s(`Statistic`,i,n),p=o(()=>{let{self:{labelFontWeight:e,valueFontSize:t,valueFontWeight:n,valuePrefixTextColor:r,labelTextColor:i,valueSuffixTextColor:a,valueTextColor:o,labelFontSize:s},common:{cubicBezierEaseInOut:l}}=c.value;return{"--n-bezier":l,"--n-label-font-size":s,"--n-label-font-weight":e,"--n-label-text-color":i,"--n-value-font-weight":n,"--n-value-font-size":t,"--n-value-prefix-text-color":r,"--n-value-suffix-text-color":a,"--n-value-text-color":o}}),m=r?t(`statistic`,void 0,p,e):void 0;return{rtlEnabled:f,mergedClsPrefix:n,cssVars:r?void 0:p,themeClass:m?.themeClass,onRender:m?.onRender}},render(){var e;let{mergedClsPrefix:t,$slots:{default:n,label:a,prefix:o,suffix:s}}=this;return(e=this.onRender)==null||e.call(this),i(`div`,{class:[`${t}-statistic`,this.themeClass,this.rtlEnabled&&`${t}-statistic--rtl`],style:this.cssVars},r(a,e=>i(`div`,{class:`${t}-statistic__label`},this.label||e)),i(`div`,{class:`${t}-statistic-value`,style:{fontVariantNumeric:this.tabularNums?`tabular-nums`:``}},r(o,e=>e&&i(`span`,{class:`${t}-statistic-value__prefix`},e)),this.value===void 0?r(n,e=>e&&i(`span`,{class:`${t}-statistic-value__content`},e)):i(`span`,{class:`${t}-statistic-value__content`},this.value),r(s,e=>e&&i(`span`,{class:`${t}-statistic-value__suffix`},e))))}});export{f as t};